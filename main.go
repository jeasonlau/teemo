package main

import (
	"flag"
	"fmt"
	"github.com/gen2brain/beeep"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"teemo/ipgw"
	"time"
)

var (
	u, p string
	f    int
	v    bool

	semesters   [][]string
	gpa, reqUrl string
	client      *http.Client

	upPath, downPath string
)

func init() {
	flag.StringVar(&u, "u", "", "学号")
	flag.StringVar(&p, "p", "", "密码")
	flag.IntVar(&f, "f", 60, "频率 单位秒")
	flag.BoolVar(&v, "v", false, "使用webvpn")
	flag.Usage = usage
}

func main() {
	flag.Parse()

	upPath, downPath = GetImgPath()

	go func() {
		// 登陆
		client = login()
		fmt.Println("登陆成功")
		var newGPA string
		for {
			newGPA = getGPA()
			fmt.Printf("%10s\t绩点: %s\n", time.Now().Format("2006-01-02 15:04:05"), newGPA)
			if newGPA != gpa {
				if newGPA == "获取失败" {
					goto wait
				}

				if len(gpa) < 1 {
					gpa = newGPA
					goto wait
				}

				n, _ := strconv.ParseFloat(newGPA, 32)
				g, _ := strconv.ParseFloat(gpa, 32)
				diff := n - g
				if diff > 0 {
					err := beeep.Notify("Teemo", "绩点变高啦", fmt.Sprintf("绩点上升了\t%.4f\n当前绩点\t%s", diff, newGPA), upPath)
					if err != nil {
						fmt.Println("推送提示失败")
					}
				} else {
					err := beeep.Notify("Teemo", "绩点降低了", fmt.Sprintf("绩点降低了\t%.4f\n当前绩点\t%s", diff, newGPA), downPath)
					if err != nil {
						fmt.Println("推送提示失败")
					}
				}
				gpa = newGPA
				//绩点改变
			}
		wait:
			time.Sleep(time.Duration(f) * time.Second)
		}
	}()

	select {}
}

func login() *http.Client {
	config := ipgw.NewLoginConfig()
	config.User = &ipgw.User{
		Username: u,
		Password: p,
	}
	config.Webvpn = v
	cookie, code := ipgw.Login(config)
	if code != 0 {
		fmt.Println(ipgw.CodeText(code))
		os.Exit(2)
	}

	c := ipgw.NewCasClient(cookie, v)
	if v {
		_, _ = c.Get("https://219-216-96-4.webvpn.neu.edu.cn/eams/homeExt.action")

		reqUrl = fmt.Sprintf("https://219-216-96-4.webvpn.neu.edu.cn/eams/teach/grade/course/person!search.action?semesterId=12&projectType=&_=%d", time.Now().Unix())
	} else {
		reqUrl = fmt.Sprintf("http://219.216.96.4/eams/teach/grade/course/person!search.action?semesterId=12&projectType=&_=%d", time.Now().Unix())
	}
	return c
}

func getGPA() string {
	resp, err := client.Get(reqUrl)
	if err != nil {
		fmt.Println("网络错误")
		return "获取失败"
	}
	body := readBody(resp)

	gpaExp := regexp.MustCompile(`<div>总平均绩点：(.+?)</div>`)
	gpa := gpaExp.FindAllStringSubmatch(body, -1)
	if len(gpa) < 1 {
		return "获取失败"
	}
	return gpa[0][1]
}

func readBody(resp *http.Response) (body string) {
	res, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return string(res)
}

func usage() {
	fmt.Println(`
teemo -u 学号 -p 密码
    如 teemo -u 2018xxxx -p abcdefg

teemo -u 学号 -p 密码 -f 监控频率(单位秒)
    如 teemo -u 2018xxxx -p abcdefg -f 60

teemo -u 学号 -p 密码 -v 使用webvpn
    如 teemo -u 2018xxxx -p abcdefg -v

若不指定u和p，默认使用ipgw保存的账号
若不指定f，默认60`)
}
