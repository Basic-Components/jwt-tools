# jwttools

jwt工具集合,主要是对jwt-go进行封装,以提供更加好用的接口

## 使用方法:

这个项目可以作为包用,也可以起一个grpc服务作为签名中心配合sdk使用.

### 当做包使用

在`utils`子模块中定义了签名器接口`Signer`和验证器接口`Verifier`.

`jwtsigner`子模块中实现了对称签名`Symmetric`(使用hash)和非对称签名`Asymmetric`(使用rsa和ecdsa)的结构体,
其中`Symmetric`使用函数`SymmetricNew`初始化;`Asymmetric`使用`AsymmetricNew`,`AsymmetricFromPEM`,`AsymmetricFromPEMFile`初始化.

这两个结构体都实现了`Signer`接口.

与`jwtsigner`子模块类似,`jwtverifier`子模块实现了对称签名校验器`Symmetric`(使用hash)和非对称签名校验器`Asymmetric`(使用rsa和ecdsa)的结构体.
其中`Symmetric`使用函数`SymmetricNew`初始化;`Asymmetric`使用`AsymmetricNew`,`AsymmetricFromPEM`,`AsymmetricFromPEMFile`初始化.
这两个结构体都实现了`Verifier`接口.

> 一个例子

```golang
package main
import (
	"github.com/Basic-Components/jwttools/jwtsigner"
	"github.com/Basic-Components/jwttools/jwtverifier"
)

func main(){
	signer, err := jwtsigner.AsymmetricFromPEMFile("RS256", "./autogen_rsa.pem")
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	jwtverifier, err := jwtverifier.AsymmetricFromPEMFile("RS256", "./autogen_rsa_pub.pem")
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	got, err := signer.SignJSONString(`{"a":1}`, "1", "a")
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	claims, err := jwtverifier.Verify(got)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	if int(claims["a"].(float64)) != 1 {
		fmt.Println("wrong")
		return
	}
	fmt.Println("ok")
}

```

### grpc服务作为签名中心

该项目的签名中心服务使用grpc作为协议,其接口定义在`schema`目录中.注意目前只实现了hs256和rs256的签名验签功能.

如果要本地安装可以使用

```bash
go get github.com/Basic-Components/jwttools/jwtcenter
go build github.com/Basic-Components/jwttools/jwtcenter
```

其启动参数为

```bash
Usage of bin/darwin-amd64/jwtrpc:
  -a, --address string            要启动的服务器地址
  -c, --config string             配置文件位置
  -k, --hash_key string           hash签名使用的盐
  -l, --loglevel string           log的等级
  -n, --name string               服务名称
  -r, --private_key_path string   签名使用的私钥地址
  -u, --public_key_path string    验签使用的公钥地址
```

也可以通过设置环境变量和设置配置文件来设置配置.默认的配置文件未知在`"/etc/jwtcenter/", "$HOME/.jwtcenter", "."`.

我们也可以使用docker来启动,镜像为`hsz1273327/jwtcenter`

## sdk

调用配置中心除了可以自己封装外也可以使用sdk,这个sdk封装了对象`RemoteCenter`实现了签名器接口`Signer`和验证器接口`Verifier`,其初始化方法是使用函数`New`来创建简单连接,也可以使用`NewWithLocalBalance`来做基于本地注册中心的负载均衡连接

```bash
go get -u -v  github.com/Basic-Components/jwttools/jwtcentersdk
```

```golang
package main
import (
	"fmt"

	jwtsdk "github.com/Basic-Components/jwttools/jwtcentersdk"
)

func main(){
	jwtcenter := jwtsdk.New("localhost:5000","RS256", time.Second)
	claims := map[string]interface{}{"IP": "127.0.0.1", "name": "S124"}
	token, err := jwtcenter.Sign(claims,"1", "a")
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	fmt.Printf("token: %v", token)
	gotclaims, err := client.Verify(token)
	if err != nil {
		fmt.Printf("got claims error: %v", err)
		return
	}
    fmt.Printf("claims: %v", gotclaims)
}
```

## change log

### v2.1.0

1. 增加了sdk的本地负载均衡方式
2. 增加了代理对象方便初始化


### v2.0.0

完成了基本结构

## todo

1. 后续会封装python版本sdk和js版本sdk
2. 后续会加上使用etcdv3做服务发现和负载均衡的方法