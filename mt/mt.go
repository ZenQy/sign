package mt

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

func sign(client *resty.Client) bool {
	url := "https://bbs.binmt.cc/"
	resp, err := client.R().Get(url)
	if err != nil {
		return false
	}
	// 与浏览器结果不一致
	ctx := resp.String()
	re := regexp.MustCompile(`<input type="hidden" name="formhash" value="(.*?)" />`)
	formhash := re.FindStringSubmatch(ctx)
	if len(formhash) != 2 {
		return false
	}
	url = "https://bbs.binmt.cc/k_misign-sign.html?operation=qiandao&format=button&inajax=1&ajaxtarget=midaben_sign&formhash=" + formhash[1]
	resp, err = client.R().Get(url)
	if err != nil {
		return false
	}
	ctx = resp.String()
	return strings.Contains(ctx, "签到成功") || strings.Contains(ctx, "今日已签")
}

// Do 签到
func Do(username, password string) string {
	client := resty.New()

	if !login(client, username, password) {
		return "登录失败"
	}

	if sign(client) {
		return "签到成功"
	}
	return "签到失败"
}

func login(client *resty.Client, username, password string) bool {
	formdata := map[string]string{
		"loginfield": "username",
		"username":   username,
		"password":   password,
		"questionid": "0",
		"cookietime": "2592000",
	}
	url := "https://bbs.binmt.cc/member.php?mod=logging&action=login"
	resp, err := client.R().Get(url)
	if err != nil {
		return false
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return false
	}
	sel := doc.Find(`form[id^="loginform_"]`)
	url, ok := sel.Attr("action")
	if !ok {
		return false
	}
	sel.Find("div > input[type=hidden]").Each(func(_ int, s *goquery.Selection) {
		name, ok1 := s.Attr("name")
		value, ok2 := s.Attr("value")
		if ok1 && ok2 && name == "formhash" {
			formdata[name] = value
		}
	})

	resp, err = client.R().SetFormData(formdata).Post("https://bbs.binmt.cc/" + url)
	if err != nil {
		return false
	}

	return strings.Contains(resp.String(), "欢迎您回来")
}
