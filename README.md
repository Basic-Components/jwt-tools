# jwttools

jwt工具集合,主要是对jwt-go进行封装,以提供更加好用的接口

## 使用方法:

这个项目可以作为包用,也可以起一个grpc服务作为签名中心配合sdk使用.

### grpc服务作为签名中心

### 当做包使用



### 服务端
```bash
go get github.com/Basic-Components/jwttools
go build github.com/Basic-Components/jwttools/jwtcenter
```

```bash
Usage of bin/darwin-amd64/jwtrpc:
  -a, --address string            要启动的服务器地址
  -c, --config string             配置文件位置
  -g, --genkey                    创建rsa公私钥对
  -i, --iss string                签名者
  -l, --loglevel string           创建rsa公私钥对 (default "WARN")
  -r, --private_key_path string   指定私钥位置
  -u, --public_key_path string    指定公钥位置
  -m, --sign_method string        指定签名方法
```

## sdk

```bash
go get -u -v  github.com/Basic-Components/jwtcentersdk
```

```golang
import (
	"fmt"

	jwtclient "github.com/Basic-Components/jwtrpc/jwtclient"
)

...
    client := jwtclient.New("localhost:5000")
	claims := map[string]interface{}{"IP": "127.0.0.1", "name": "S124"}
	token, err := client.GetToken(claims)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	fmt.Printf("token: %v", token)
	gotclaims, err := client.VerifyToken(token)
	if err != nil {
		fmt.Printf("got claims error: %v", err)
		return
	}
    fmt.Printf("claims: %v", gotclaims)
...
```

