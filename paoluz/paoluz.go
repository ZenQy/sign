package paoluz

import (
	"encoding/json"
	"strings"

	"github.com/go-resty/resty/v2"
)

type response struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

// Do 签到
func Do(username, password string) string {
	client := resty.New()

	if !login(client, username, password) {
		return "登录失败"
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

func login(client *resty.Client, username, password string) bool {
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
