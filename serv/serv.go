package serv

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	sync "sync"
	"syscall"
	"time"

	log "github.com/Golang-Tools/loggerhelper"

	"github.com/liyue201/grpc-lb/common"
	"github.com/liyue201/grpc-lb/registry"
	zk "github.com/liyue201/grpc-lb/registry/zookeeper"

	"github.com/Basic-Components/jwttools/idgener"
	"github.com/Basic-Components/jwttools/interfaces"
	declare "github.com/Basic-Components/jwttools/jwtrpcdeclare"
	"github.com/Basic-Components/jwttools/jwtsigner"
	"github.com/Basic-Components/jwttools/jwtverifier"
	"github.com/Basic-Components/jwttools/options"
	"github.com/Basic-Components/jwttools/utils"
	se "github.com/Golang-Tools/schema-entry-go"
	grpc "google.golang.org/grpc"
	channelz "google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

//Server grpc的服务器结构体
//服务集成了如下特性:
//设置收发最大消息长度
//健康检测
//gzip做消息压缩
//接口反射
//channelz支持
//TLS支持
//keep alive 支持
type Server struct {
	App_Name       string   `json:"app_name,omitempty" jsonschema:"description=服务名"`
	App_Version    string   `json:"app_version,omitempty" jsonschema:"description=服务版本"`
	Address        string   `json:"address,omitempty" jsonschema:"description=服务的主机和端口"`
	Log_Level      string   `json:"log_level,omitempty" jsonschema:"description=项目的log等级,enum=TRACE,enum=DEBUG,enum=INFO,enum=WARN,enum=ERROR"`
	Zookeeper_URL  []string `json:"zookeeper_url,omitempty" jsonschema:"description=负载均衡使用的zookeeper地址序列以逗号分隔"`
	Balance_Weight string   `json:"balance_weight,omitempty" jsonschema:"description=负载均衡的权重"`

	Max_Recv_Msg_Size int `json:"max_rec_msg_size,omitempty" jsonschema:"description=允许接收的最大消息长度"`
	Max_Send_Msg_Size int `json:"max_send_msg_size,omitempty" jsonschema:"description=允许发送的最大消息长度"`

	Initial_Window_Size      int `json:"initial_window_size,omitempty" jsonschema:"description=基于Stream的滑动窗口大小"`
	Initial_Conn_Window_Size int `json:"initial_conn_window_size,omitempty" jsonschema:"description=基于Connection的滑动窗口大小"`
	Max_Concurrent_Streams   int `json:"max_concurrent_streams,omitempty" jsonschema:"description=一个连接中最大并发Stream数"`

	Use_Channelz                                bool   `json:"use_channelz,omitempty" jsonschema:"description=是否使用channelz协助优化"`
	Server_Cert_Path                            string `json:"server_cert_path ,omitempty" jsonschema:"description=使用TLS时服务端的证书位置"`
	Server_Key_Path                             string `json:"server_key_path,omitempty" jsonschema:"description=使用TLS时服务端证书的私钥位置"`
	Max_Connection_Idle                         int    `json:"max_connection_idle,omitempty" jsonschema:"description=客户端连接的最大空闲时长"`
	Max_Connection_Age                          int    `json:"max_connection_age,omitempty" jsonschema:"description=如果连接存活超过n则发送goaway"`
	Max_Connection_Age_Grace                    int    `json:"max_connection_age_grace,omitempty" jsonschema:"description=强制关闭连接之前允许等待的rpc在n秒内完成"`
	Keepalive_Time                              int    `json:"keepalive_time,omitempty" jsonschema:"description=空闲连接每隔n秒ping一次客户端已确保连接存活"`
	Keepalive_Timeout                           int    `json:"keepalive_timeout,omitempty" jsonschema:"description=ping时长超过n则认为连接已死"`
	Keepalive_Enforcement_Min_Time              int    `json:"keepalive_enforement_min_time,omitempty" jsonschema:"description=如果客户端超过每n秒ping一次则终止连接"`
	Keepalive_Enforcement_Permit_Without_Stream bool   `json:"keepalive_enforement_permit_without_stream,omitempty" jsonschema:"description=即使没有活动流也允许ping"`

	Algo_Type                  string `json:"algo_type" jsonschema:"required,description=服务使用的加密算法,enum=RS256,enum=RS384,enum=RS512,enum=ES256,enum=ES384,enum=ES512,enum=HS256,enum=HS384,enum=HS512""`
	Secret_Key                 string `json:"secret_key,omitempty" jsonschema:"description=对称加密构造hash的盐"`
	Private_Key_Path           string `json:"secret_key,omitempty" jsonschema:"description=非对称加密的私钥"`
	Public_Key_Path            string `json:"public_key_path,omitempty" jsonschema:"description=非对称加密的公钥"`
	Default_TTL                int    `json:"secret_key,omitempty" jsonschema:"description=jwt默认的过期时长,单位s"`
	Default_Effective_Interval int    `json:"secret_key,omitempty" jsonschema:"description=jwt默认的生效间隔,单位s"`
	Jti_Gen                    string `json:"jti_gen,omitempty" jsonschema:"description=生成jti的生成器,enum=uuid4,enum=sonyflake"`

	service       *registry.ServiceInfo
	healthservice *health.Server
	registrar     *zk.Registrar

	signer   interfaces.CanSign
	verifier interfaces.CanVerify
}

//Main 服务的入口函数
func (s *Server) Main() {
	// 初始化log
	log.Init(s.Log_Level, log.Dict{
		"app_name":    s.App_Name,
		"app_version": s.App_Version,
	})
	log.Info("获得参数", log.Dict{"ServiceConfig": s})
	opts := []options.SignerOption{options.WithIss(s.App_Name)}
	if s.Default_TTL > 0 {
		opts = append(opts, options.WithDefaultTTL(time.Duration(s.Default_TTL)*time.Second))
	}
	if s.Default_Effective_Interval > 0 {
		opts = append(opts, options.WithDefaultEffectiveInterval(time.Duration(s.Default_Effective_Interval)*time.Second))
	}
	if s.Jti_Gen != "" {
		switch s.Jti_Gen {
		case "uuid4":
			{
				opts = append(opts, options.WithJtiGen(&idgener.UUID4Gen{}))
			}
		case "sonyflake":
			{
				opts = append(opts, options.WithJtiGen(idgener.NewSonyflakeGen()))
			}
		}
	}
	if s.Public_Key_Path != "" {

	} else {
		algo, err := utils.AlgoStrTOAlgoEnum(s.Algo_Type)
		if err != nil {
			log.Error("AlgoStrTOAlgoEnum error")
			os.Exit(2)
		}
		if utils.IsAsymmetric(algo) {
			if s.Private_Key_Path == "" || s.Public_Key_Path == "" {
				log.Error("Asymmetric algo must set Private_Key_Path and Public_Key_Path")
				os.Exit(2)
			}
			signer, err := jwtsigner.AsymmetricFromPEMFile(algo, s.Private_Key_Path, opts...)
			if err != nil {
				log.Error("jwtsigner.SymmetricNew error", log.Dict{"err": err, "algo": algo.String(), "Private_Key_Path": s.Private_Key_Path})
				os.Exit(2)
			}
			s.signer = signer
			verifier, err := jwtverifier.AsymmetricFromPEMFile(algo, s.Public_Key_Path)
			if err != nil {
				log.Error("jwtverifier.SymmetricNew", log.Dict{"err": err, "algo": algo.String(), "Public_Key_Path": s.Public_Key_Path})
				os.Exit(2)
			}
			s.verifier = verifier

		} else {
			if s.Secret_Key == "" {
				log.Error("Symmetric algo must set secret_key")
				os.Exit(2)
			}
			signer, err := jwtsigner.SymmetricNew(algo, s.Secret_Key, opts...)
			if err != nil {
				log.Error("jwtsigner.SymmetricNew error", log.Dict{"err": err, "algo": algo.String()})
				os.Exit(2)
			}
			s.signer = signer
			verifier, err := jwtverifier.SymmetricNew(algo, s.Secret_Key)
			if err != nil {
				log.Error("jwtverifier.SymmetricNew", log.Dict{"err": err, "algo": algo.String()})
				os.Exit(2)
			}
			s.verifier = verifier
		}

	}
	s.Run()
}

//RunServer 启动服务
func (s *Server) RunServer() {
	lis, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Error("Failed to Listen", log.Dict{"error": err.Error(), "address": s.Address})
		os.Exit(1)
	}
	opts := []grpc.ServerOption{}
	if s.Max_Recv_Msg_Size != 0 {
		opts = append(opts, grpc.MaxRecvMsgSize(s.Max_Recv_Msg_Size))
	}
	if s.Max_Send_Msg_Size != 0 {
		opts = append(opts, grpc.MaxSendMsgSize(s.Max_Send_Msg_Size))
	}
	if s.Initial_Window_Size != 0 {
		opts = append(opts, grpc.InitialWindowSize(int32(s.Initial_Window_Size)))
	}
	if s.Initial_Conn_Window_Size != 0 {
		opts = append(opts, grpc.InitialConnWindowSize(int32(s.Initial_Conn_Window_Size)))
	}
	if s.Max_Concurrent_Streams != 0 {
		opts = append(opts, grpc.MaxConcurrentStreams(uint32(s.Max_Concurrent_Streams)))
	}
	if s.Server_Cert_Path != "" && s.Server_Key_Path != "" {
		creds, err := credentials.NewServerTLSFromFile(s.Server_Cert_Path, s.Server_Key_Path)
		if err != nil {
			log.Warn("Failed to Listen as a TLS Server", log.Dict{"error": err.Error()})
		}
		opts = append(opts, grpc.Creds(creds))
	}
	if s.Max_Connection_Idle != 0 || s.Max_Connection_Age != 0 || s.Max_Connection_Age_Grace != 0 || s.Keepalive_Time != 0 || s.Keepalive_Timeout != 0 {
		kasp := keepalive.ServerParameters{
			MaxConnectionIdle:     time.Duration(s.Max_Connection_Idle) * time.Second,
			MaxConnectionAge:      time.Duration(s.Max_Connection_Age) * time.Second,
			MaxConnectionAgeGrace: time.Duration(s.Max_Connection_Age_Grace) * time.Second,
			Time:                  time.Duration(s.Keepalive_Time) * time.Second,
			Timeout:               time.Duration(s.Keepalive_Timeout) * time.Second,
		}
		opts = append(opts, grpc.KeepaliveParams(kasp))
	}

	if s.Keepalive_Enforcement_Min_Time != 0 || s.Keepalive_Enforcement_Permit_Without_Stream == true {
		kaep := keepalive.EnforcementPolicy{
			MinTime:             time.Duration(s.Keepalive_Enforcement_Min_Time) * time.Second,
			PermitWithoutStream: s.Keepalive_Enforcement_Permit_Without_Stream,
		}
		opts = append(opts, grpc.KeepaliveEnforcementPolicy(kaep))
	}
	gs := grpc.NewServer(opts...)
	defer gs.Stop()
	s.healthservice = health.NewServer()
	healthpb.RegisterHealthServer(gs, s.healthservice)

	declare.RegisterJwtServiceServer(gs, s)
	reflection.Register(gs)
	if s.Use_Channelz {
		channelz.RegisterChannelzServiceToServer(gs)
	}
	log.Info("Server Start", log.Dict{"address": s.Address})
	err = gs.Serve(lis)
	if err != nil {
		log.Error("Failed to Serve", log.Dict{"error": err})
		os.Exit(1)
	}
}

//RegistService 注册服务到zookeeper
func (s *Server) RegistService() {
	if s.registrar != nil && s.service != nil {
		log.Warn("服务注册已经初始化")
		return
	}
	port := strings.Split(s.Address, ":")[1]
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Error("获取本地ip失败", log.Dict{"place": "RegistService", "err": err})
		os.Exit(1)
	}
	ip := ""
	for _, _ip := range addrs {
		IP := _ip.String()
		if strings.HasPrefix(IP, "172.16.1.") {
			if strings.Contains(IP, "/") {
				ip = strings.Split(IP, "/")[0]
			} else {
				ip = IP
			}
			break
		}
	}
	if ip == "" {
		log.Error("未找到ip", log.Dict{"place": "RegistService"})
		os.Exit(1)
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Error("获取本地容器hostname失败", log.Dict{"place": "RegistService", "err": err})
		os.Exit(1)
	}
	service := &registry.ServiceInfo{
		InstanceId: hostname,
		Name:       s.App_Name,
		Version:    s.App_Version,
		Address:    fmt.Sprintf("%s:%s", ip, port),
		Metadata:   metadata.Pairs(common.WeightKey, s.Balance_Weight),
	}
	log.Info("注册的服务", log.Dict{"service": *service})
	registrar, err := zk.NewRegistrar(
		&zk.Config{
			ZkServers:      s.Zookeeper_URL,
			RegistryDir:    "/backend/services",
			SessionTimeout: time.Second,
		})
	if err != nil {
		log.Error("regist error", log.Dict{"err": err})
		os.Exit(1)
	}
	s.registrar = registrar
	s.service = service

}

//Run 执行grpc服务
func (s *Server) Run() {
	if len(s.Zookeeper_URL) == 0 || s.Zookeeper_URL == nil {
		s.RunServer()
	} else {
		s.RegistService()
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			s.RunServer()
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			s.registrar.Register(s.service)
			wg.Done()
		}()
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
		s.registrar.Unregister(s.service)
		// serv.Stop()
		wg.Wait()
	}
}

var Endpoint, _ = se.New(&se.EntryPointMeta{Name: "jwtcenter", Usage: "jwtcenter [options]"}, &Server{
	App_Name:                   "jwtcenter",
	App_Version:                "3.0.0",
	Address:                    "0.0.0.0:5000",
	Log_Level:                  "INFO",
	Algo_Type:                  "HS256",
	Secret_Key:                 "guesssecrest",
	Default_TTL:                7 * 24 * 3600,
	Default_Effective_Interval: 180,
	Jti_Gen:                    "uuid4",
})
