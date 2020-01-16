package ipgw

import "net/http"

type User struct {
	Username string
	Password string
	Cookie   string
}

type LoginConfig struct {
	User   *User
	Webvpn bool
}

type ProxyConfig struct {
	User       *User
	LaunchUrl  string
	ServiceUrl string
	Method     string
	Headers    *http.Header
	Body       string
}

func NewLoginConfig() *LoginConfig {
	return &LoginConfig{User: &User{}}
}

func NewProxyConfig() *ProxyConfig {
	return &ProxyConfig{User: &User{}, Headers: &http.Header{}, Method: "GET"}
}
