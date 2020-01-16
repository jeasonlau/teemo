<p align="center">
    <img src="https://raw.githubusercontent.com/iMyOwn/teemo/master/img/logo.png" alt="logo" width="200">
</p>

<h1 align="center">Teemo</h1>
<p align="center">
    <img src="https://img.shields.io/github/release-date/iMyOwn/teemo" alt="">
    <img src="https://img.shields.io/github/license/iMyown/teemo" alt="">
    <img src="https://img.shields.io/github/go-mod/go-version/iMyOwn/teemo" alt="">
</p>


> 东北大学GPA监控程序

# 修改内容

由于个人使用的桌面环境为`xfce`，无法正常接收通知，故`fork`该项目并将其通知方式修改为邮件。

# 使用方法

以QQ邮箱为例：
1. 登录官网获取`smtp`授权码（此处不赘述）。

2. 打开`mail/mail.go`，将8/9行修改为：

   ```go
   	from := "你的邮箱"
   	password := "你的授权码"
   ```

3. 打开`main.go`，将59/64行修改为：

   ```go
   err := mail.SendMail("你的邮箱", fmt.Sprintf("绩点(升高|降低)了\t%.4f\n当前绩点\t%s", diff, newGPA))
   ```

4. 在项目根目录执行：

   ```bash
   make build-linux
   //or
   make build-windows
   ```
   位于`build/(linux|windows)`目录的即为可执行文件。
   
# 系统要求
满足以下之一即可
- Windows 8+
- Linux

# 下载
`teemo`现已成为`ipgw`的一个工具，可通过`ipgw tool get`获取，详见[IPGW Tool](https://github.com/iMyOwn/ipgw)

配置`ipgw`妥当后，可快速使用`teemo`
```shell script
ipgw tool get teemo
```

查看使用说明
```shell script
teemo -h
```

若正常输出使用说明，则已下载完毕，否则请尝试再次`ipgw tool get teemo`重新下载

# 更新

使用`ipgw`的工具更新功能即可
```shell script
ipgw tool update teemo
```

在特殊情况下可能需要强制更新
```shell script
ipgw tool update -f teemo
```

# 使用
本工具为命令行工具，Linux请在`terminal`中使用，windows请在`cmd`中使用。

自`v1.1.0`起已支持外网通过webvpn监控。

```
teemo -u 学号 -p 密码
    如 teemo -u 2018xxxx -p abcdefg

teemo -u 学号 -p 密码 -f 监控频率(单位秒)
    如 teemo -u 2018xxxx -p abcdefg -f 60

teemo -u 学号 -p 密码 -v 使用webvpn
    如 teemo -u 2018xxxx -p abcdefg -v

若不指定u和p，默认使用ipgw保存的账号
若不指定f，默认60
```

> 监控中请不要关闭程序
> 
> - Linux不要关闭Terminal; 也可使用类似于`screen`、`nohup`等工具与命令后台运行程序
> - Windows不要关闭CMD

# 修改程序
【注意】编写过程中修改了部分依赖包代码以解决依赖包的Bug，因此请使用`/vendor`中提供的依赖包

编译时添加flag `-mod=vendor`

# 开源协议
Mit License.