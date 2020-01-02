package main

import (
	"flag"
	"fmt"
	"github.com/gen2brain/beeep"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	jar       http.CookieJar
	client    *http.Client
	u, p      string
	f, s      int
	q         bool
	semesters [][]string
	gpa       string
)

func init() {
	jar, _ = cookiejar.New(nil)
	client = &http.Client{Jar: jar}

	flag.StringVar(&u, "u", "", "学号")
	flag.StringVar(&p, "p", "", "密码")
	flag.IntVar(&s, "s", 12, "学期代码")
	flag.IntVar(&f, "f", 60, "频率 单位秒")
	flag.BoolVar(&q, "q", false, "查询学期代码")
	flag.Usage = usage
}

func main() {
	flag.Parse()
	login()
	if q {
		// 获取学期信息
		getSemesters()
		return
	}

	go func() {
		for {
			newGPA := getGPA()
			fmt.Printf("%10s\t绩点: %s\n", time.Now().Format("2006-01-02 15:04:05"), newGPA)
			if newGPA != gpa {
				n, _ := strconv.ParseFloat(newGPA, 32)
				g, _ := strconv.ParseFloat(gpa, 32)
				diff := n - g
				if diff > 0 {
					err := beeep.Notify("Teemo", "绩点变高啦", fmt.Sprintf("绩点上升了%f", diff), "img/up.png")
					if err != nil {
						fmt.Println("推送提示失败")
					}
				} else {
					err := beeep.Notify("Teemo", "绩点降低了", fmt.Sprintf("绩点降低了%f", diff), "img/down.png")
					if err != nil {
						fmt.Println("推送提示失败")
					}
				}
				gpa = newGPA
				//绩点改变
			}
			time.Sleep(time.Duration(s) * time.Second)
		}
	}()

	select {}
}

func login() {
	reqUrl := "https://pass.neu.edu.cn/tpass/login?service=http%3A%2F%2F219.216.96.4%2Feams%2FhomeExt.action"
	// 请求获得必要参数
	resp, err := client.Get(reqUrl)
	handlerErr(err)

	// 读取响应内容
	body := readBody(resp)

	// 读取lt
	ltExp := regexp.MustCompile(`name="lt" value="(.+?)"`)
	lts := ltExp.FindAllStringSubmatch(body, -1)

	postUrlExp := regexp.MustCompile(`id="loginForm" action="(.+?)"`)
	postUrls := postUrlExp.FindAllStringSubmatch(body, -1)

	if len(lts) < 1 || len(postUrls) < 1 {
		fmt.Println("参数获取失败")
		os.Exit(2)
	}
	lt := lts[0][1]
	postUrl := postUrls[0][1]

	// 拼接data
	data := "rsa=" + u + p + lt +
		"&ul=" + strconv.Itoa(len(u)) +
		"&pl=" + strconv.Itoa(len(p)) +
		"&lt=" + lt +
		"&execution=e1s1" +
		"&_eventId=submit"

	// 构造请求
	req, _ := http.NewRequest("POST",
		"https://pass.neu.edu.cn"+postUrl,
		strings.NewReader(data))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "pass.neu.edu.cn")
	req.Header.Add("Origin", "https://pass.neu.edu.cn")
	req.Header.Add("Referer", reqUrl)

	// 发送请求
	resp, err = client.Do(req)
	handlerErr(err)

	// 读取响应内容
	body = readBody(resp)

	idExp := regexp.MustCompile(`class="personal-name"> .+?\((\d+?)\) </a>`)
	ids := idExp.FindAllStringSubmatch(body, -1)

	if len(ids) < 1 {
		fmt.Println("登陆失败")
		os.Exit(2)
	}
}

func getSemesters() {
	data := "tagId=semesterBar1111111111Semester&dataType=semesterCalendar&value=49&empty=false"
	req, _ := http.NewRequest("POST", "http://219.216.96.4/eams/dataQuery.action", strings.NewReader(data))
	req.Header.Set("Referer", "http://219.216.96.4/eams/homeExt.action")
	req.Header.Set("Origin", "http://219.216.96.4")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := client.Do(req)
	handlerErr(err)
	body := readBody(resp)

	semExp := regexp.MustCompile(`{id:(\d+?),schoolYear:"(.+?)",name:"(.+?)"}`)
	semesters = semExp.FindAllStringSubmatch(body, -1)
	// 学期太多啦，忽略掉前三十个
	for _, semester := range semesters[30:] {
		fmt.Printf("%s学年%s学期\t%s\n", semester[2], semester[3], semester[1])
	}
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
	fmt.Println(`监控绩点	teemo -u 学号 -p 密码 -s 学期代码
指定监控频率	teemo -u 学号 -p 密码 -s 学期代码 -f 监控频率(单位秒)
查询学期代码	teemo -u 学号 -p 密码 -q`)
}
