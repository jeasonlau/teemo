package ipgw

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// 登陆
func Login(config *LoginConfig) (cookie string, code int) {
	if config.Webvpn {
		return execCommand("ipgw", "api", "v1", "login", "-u", config.User.Username, "-p", config.User.Password, "-c", config.User.Cookie, "-v")
	}
	return execCommand("ipgw", "api", "v1", "login", "-u", config.User.Username, "-p", config.User.Password, "-c", config.User.Cookie)
}

// 代理请求
func Proxy(config *ProxyConfig) (body string, code int) {
	headers, _ := json.Marshal(config.Headers)
	return execCommand("ipgw", "api", "v1", "proxy",
		"-u", config.User.Username, "-p", config.User.Password, "-c", config.User.Cookie,
		"-s", config.ServiceUrl,
		"-m", config.Method,
		"-h", string(headers),
		"-b", config.Body)
}

// 创建客户端
func NewCasClient(cookie string, webvpn bool) (client *http.Client) {
	n := &http.Client{Timeout: 6 * time.Second}
	jar, _ := cookiejar.New(nil)
	// 绑定session
	n.Jar = jar
	if webvpn {
		jar.SetCookies(&url.URL{
			Scheme: "https",
			Host:   "pass-443.webvpn.neu.edu.cn",
			Path:   "/tpass/",
		}, []*http.Cookie{
			{
				Name:   "CASTGC",
				Value:  cookie,
				Domain: "pass-443.webvpn.neu.edu.cn",
				Path:   "/tpass/",
			},
		})
	} else {
		jar.SetCookies(&url.URL{
			Scheme: "https",
			Host:   "pass.neu.edu.cn",
			Path:   "/tpass/",
		}, []*http.Cookie{
			{
				Name:   "CASTGC",
				Value:  cookie,
				Domain: "pass.neu.edu.cn",
				Path:   "/tpass/",
			},
		})
	}
	return n
}

// 当调用命令失败时返回-1。
// 当命令返回错误码非0时返回空字符串与错误码。
// 否则返回调用结果与0.
func execCommand(name string, params ...string) (result string, code int) {
	var outbuf, errbuf bytes.Buffer

	cmd := exec.Command(name, params...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		return "", -1
	}
	result = string(outbuf.Bytes())
	c := errbuf.Bytes()
	c = c[:len(c)-1]
	code, _ = strconv.Atoi(string(c))
	return
}
