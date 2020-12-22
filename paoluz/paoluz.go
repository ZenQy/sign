package paoluz

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
)

type response struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

// Sign 签到
func Sign(m map[string]string) string {
	username, password, cookie := m["username"], m["password"], m["cookie"]
	client := resty.New()
	if !loginCookie(client, cookie) {
		if !loginPassword(client, username, password) {
			return "失败"
		}
		cookies := make([]string, 0)
		u, err := url.Parse("https://paoluz.com/")
		if err != nil {
			return "失败"
		}

		for _, c := range client.GetClient().Jar.Cookies(u) {
			cookies = append(cookies, c.Name+"="+c.Value)
		}
		m["cookie"] = strings.Join(cookies, "; ")
	}

	resp, err := client.R().Post("https://paoluz.com/user/checkin")
	if err != nil {
		return "失败"
	}
	res := response{}
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return "失败"
	}
	if strings.Contains(res.Msg, "流量") || strings.Contains(res.Msg, "签到") {
		return "成功"
	}
	return "失败"
}

func loginPassword(client *resty.Client, username, password string) bool {
	formdata := map[string]string{
		"email":       username,
		"passwd":      password,
		"remember_me": "on",
	}
	res := response{}
	resp, err := client.R().SetFormData(formdata).Post("https://paoluz.com/auth/login")
	if err != nil {
		return false
	}
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return false
	}
	return res.Msg == "登录成功"
}

func loginCookie(client *resty.Client, cookie string) bool {
	req := client.R()
	resp, err := req.SetHeader("cookie", cookie).Get("https://paoluz.com/user")
	if err != nil {
		return false
	}
	if strings.Contains(resp.String(), "钱包余额") {
		client.SetCookies(req.Cookies)
	}
	return false
}
