package tools

import (
	"bytes"
	"log"
	"os"
)

func Out(mod string, msg string, flag bool) {
	if flag {
		log.Printf("【%s】 : %s", mod, msg)
	}
}
func CheckClassname(path string, clname string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic("类名读取失败，请确认指定文件存在且有权限")
	}
	if !bytes.Contains(data, []byte(clname)) {
		panic("类名和文件名不一致，文件名请用 类名.xxx 的形式(xxx为任意字符串)，如:Evil.class")
	}
}

//提取class文件中的类名
