package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// IsEs 判断文件是不是ES方法加密
func IsEs(method string) bool {
	return strings.HasPrefix(method, "ES")
}

// IsRs 判断文件是不是RS方法加密
func IsRs(method string) bool {
	return strings.HasPrefix(method, "RS") || strings.HasPrefix(method, "PS")
}

// LoadData 读取并加载文件数据
func LoadData(p string) ([]byte, error) {
	if p == "" {
		return nil, fmt.Errorf("No path specified")
	}

	var rdr io.Reader
	if p == "-" {
		rdr = os.Stdin
	} else if p == "+" {
		return []byte("{}"), nil
	} else {
		if f, err := os.Open(p); err == nil {
			rdr = f
			defer f.Close()
		} else {
			return nil, err
		}
	}
	return ioutil.ReadAll(rdr)
}
