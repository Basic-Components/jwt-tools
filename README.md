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
		
		return
	}
	jwtverifier, err := jwtverifier.AsymmetricFromPEMFile("RS256", "./autogen_rsa_pub.pem")
	if err != nil {
		return
	}
	got, err := signer.SignJSONString(`{"a":1}`, "1", "a")
	if err != nil {
		
		return
	}
	claims, err := jwtverifier.Verify(got)
	if err != nil {
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


该项目的签名中心服务使用grpc作为协议,

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

