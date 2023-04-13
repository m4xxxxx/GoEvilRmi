package tools

import (
	"log"
	"os"
)

func Out(mod string, msg string, flag bool) {
	if flag {
		log.Printf("【%s】 : %s", mod, msg)
	}
}
func Getclassname(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic("类名读取失败，请确认指定文件存在且有权限")
	}
	CNL := int64(data[12])
	return string(data[13 : 13+CNL])
}

//提取class文件中的类名
