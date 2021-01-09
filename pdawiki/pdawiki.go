package pdawiki

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Sign 签到
func Sign(m map[string]string) string {
	username, password, cookie := m["username"], m["password"], m["cookie"]
	client := resty.New()
	if !loginCookie(client, cookie) {
		if !loginPassword(client, username, password) {
			return "失败"
		}
		cookies := make([]string, 0)
		u, err := url.Parse("https://www.pdawiki.com/")
		if err != nil {
			return "失败"
		}

		for _, c := range client.GetClient().Jar.Cookies(u) {
			cookies = append(cookies, c.Name+"="+c.Value)
		}
		m["cookie"] = strings.Join(cookies, "; ")

		resp, err := client.R().Get("https://www.pdawiki.com/forum/home.php?mod=spacecp&ac=profile")
		if err != nil {
			return "失败"
		}
		if !strings.Contains(resp.String(), "基本资料") {
			return "失败"
		}
	}

	return "成功"
}

func loginPassword(client *resty.Client, username, password string) bool {
	hash := md5.Sum([]byte(password))
	passwordMD5 := fmt.Sprintf("%x", hash)

	formdata := map[string]string{
		"fastloginfield": "username",
		"username":       username,
		"cookietime":     "2592000",
		"password":       passwordMD5,
		"quickforward":   "yes",
		"handlekey":      "ls",
	}
	_, err := client.R().SetFormData(formdata).Post("https://www.pdawiki.com/forum/member.php?mod=logging&action=login&loginsubmit=yes&infloat=yes&lssubmit=yes&inajax=1")
	if err != nil {
		return false
	}

	return true
}

func loginCookie(client *resty.Client, cookie string) bool {
	req := client.R()
	resp, err := req.SetHeader("cookie", cookie).Get("https://www.pdawiki.com/forum/home.php?mod=spacecp&ac=profile")
	if err != nil {
		return false
	}
	if strings.Contains(resp.String(), "基本资料") {
		client.SetCookies(req.Cookies)
	}
	return false
}
