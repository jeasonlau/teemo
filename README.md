<p align="center">
    <img src="https://raw.githubusercontent.com/iMyOwn/teemo/master/img/logo.png" alt="logo" width="200">
</p>

# Teemo
> 东北大学GPA监控程序

<p align="center">
    <img src="https://raw.githubusercontent.com/iMyOwn/teemo/master/img/demo@windows.png" alt="windows demo">
</p>
<p align="center">
    <img src="https://raw.githubusercontent.com/iMyOwn/teemo/master/img/demo@linux.png" alt="linux demo">
</p>

# 系统要求
满足以下之一即可
- OSX
- Windows 8+
- Linux

# 下载
1. 进入[Release页面](https://github.com/iMyOwn/teemo/releases)下载对应系统的压缩包
2. 解压压缩包到某个单独文件夹内
3. OSX/Linux使用`terminal`，Windows使用`cmd`，输入`teemo -h`，若输出帮助信息则完成下载
> OSX/Linux系统可能需要通过`chmod +x teemo`赋予程序可执行权限

# 使用
```shell script
# 监控绩点	
teemo -u 学号 -p 密码 -s 学期代码

# 指定监控频率
teemo -u 学号 -p 密码 -s 学期代码 -f 监控频率(单位秒)

# 查询学期代码	
teemo -u 学号 -p 密码 -q
```
> 监控中请不要关闭程序
> 
> - Linux / OSX不要关闭Terminal，或使用类似于`screen`、`nohup`等工具与命令后台运行程序
> - Windows不要关闭CMD)

# 修改程序
【注意】编写过程中修改了部分依赖包代码以解决依赖包的Bug，因此请使用`/vendor`中提供的依赖包