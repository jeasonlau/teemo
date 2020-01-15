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
	u, p      string
	f, s      int
	q         bool
	semesters [][]string
	gpa       string
	client    *http.Client

	upPath   string
	downPath string
)

func init() {
	flag.StringVar(&u, "u", "", "学号")
	flag.StringVar(&p, "p", "", "密码")
	flag.IntVar(&s, "s", 12, "学期代码")
	flag.IntVar(&f, "f", 60, "频率 单位秒")
	flag.BoolVar(&q, "q", false, "查询学期代码")
	flag.Usage = usage
}

func main() {
	flag.Parse()
	if q {
		// 获取学期信息
		fmt.Println("正在获取学期信息...")
		getSemesters()
		return
	}

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
					continue
				}

				if len(gpa) < 1 {
					gpa = newGPA
					continue
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
			time.Sleep(time.Duration(f) * time.Second)
		}
	}()

	select {}
}

func getSemesters() {
	config := ipgw.NewProxyConfig()
	config.Method = "POST"
	config.Body = "tagId=semesterBar19319741991Semester&dataType=semesterCalendar&value=12&empty=false"
	config.Headers.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	config.User = &ipgw.User{
		Username: u,
		Password: p,
	}
	config.ServiceUrl = "http://219.216.96.4/eams/dataQuery.action"
	config.Headers.Set("Referer", "http://219.216.96.4/eams/homeExt.action")
	config.Headers.Set("Origin", "http://219.216.96.4")

	body, code := ipgw.Proxy(config)
	if code != 0 {
		fmt.Println(ipgw.CodeText(code))
		os.Exit(2)
	}
	semExp := regexp.MustCompile(`{id:(\d+?),schoolYear:"(.+?)",name:"(.+?)"}`)
	semesters = semExp.FindAllStringSubmatch(body, -1)
	for _, semester := range semesters {
		fmt.Printf("%s学年%s学期\t%s\n", semester[2], semester[3], semester[1])
	}
}

func login() *http.Client {
	config := ipgw.NewLoginConfig()
	config.User = &ipgw.User{
		Username: u,
		Password: p,
	}
	cookie, code := ipgw.Login(config)
	if code != 0 {
		fmt.Println(ipgw.CodeText(code))
		os.Exit(2)
	}

	return ipgw.NewCasClient(cookie, false)
}

func getGPA() string {
	resp, err := client.Get(fmt.Sprintf("http://219.216.96.4/eams/teach/grade/course/person!search.action?semesterId=%d&projectType=&_=%d", s, time.Now().Unix()))
	handlerErr(err)
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

func handlerErr(err error) {
	if err != nil {
		fmt.Println("网络错误")
		os.Exit(2)
	}
}

func usage() {
	fmt.Println(`
监控绩点:
	teemo -u 学号 -p 密码 -s 学期代码
	teemo -u 学号 -p 密码 -s 学期代码 -f 监控频率(单位秒)
	teemo -s 学期代码	(使用ipgw保存的账号)
查询学期代码:
	teemo -u 学号 -p 密码 -q
	teemo -q		(使用ipgw保存的账号)`)
}
