package pdawiki

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

// Do 签到
func Do(username, password string) string {
	client := resty.New()

	if !login(client, username, password) {
		return "登录失败"
	}

	return extcredit(client)
}

func login(client *resty.Client, username, password string) bool {
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

func extcredit(client *resty.Client) string {
	msg := "获取积分失败"
	url := "https://www.pdawiki.com/forum/home.php?mod=spacecp&ac=credit&showcredit=1"
	resp, err := client.R().Get(url)
	if err != nil {
		log.Println(err)
		return msg
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		log.Println(err)
		return msg
	}
	txt := doc.Find(`#extcreditmenu`).Text()
	if txt == "" {
		return msg
	}
	return txt
}
