package cloud189

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-resty/resty/v2"
)

func sign(client *resty.Client) int {
	rand := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	url := "https://api.cloud.189.cn/mkt/userSign.action?rand=" + rand + "&clientType=TELEANDROID&version=8.6.3&model=SM-G930K"
	headers := map[string]string{
		"Referer":         "https://m.cloud.189.cn/zhuanti/2016/sign/index.jsp?albumBackupOpened=1",
		"Host":            "m.cloud.189.cn",
		"Accept-Encoding": "gzip, deflate",
	}

	resp, err := client.R().SetHeaders(headers).Get(url)
	if err != nil {
		log.Println(err)
		return 0
	}
	netdiskBonus, err := jsonparser.GetInt(resp.Body(), "netdiskBonus")
	if err != nil {
		log.Println(err)
		return 0
	}
	return int(netdiskBonus)
}

func login(client *resty.Client, username, password string) bool {
	url := "https://cloud.189.cn/udb/udb_login.jsp?pageId=1&redirectURL=/main.action"
	resp, err := client.R().Get(url)
	if err != nil {
		log.Println(err)
		return false
	}
	ctx := resp.String()

	captchaToken := regexpString(`captchaToken' value='(.+?)'`, ctx)

	lt := regexpString(`var lt = "(.+?)";`, ctx)
	returnUrl := regexpString(`returnUrl = '(.+?)',`, ctx)
	paramId := regexpString(`var paramId = "(.+?)";`, ctx)
	jRsaKey := regexpString(`id="j_rsaKey" value="(.+?)"`, ctx)
	key := "-----BEGIN PUBLIC KEY-----\n" + jRsaKey + "\n-----END PUBLIC KEY-----"
	url = "https://open.e.189.cn/api/logbox/oauth2/loginSubmit.do"
	headers := map[string]string{
		"lt":         lt,
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:74.0) Gecko/20100101 Firefox/76.0",
		"Referer":    "https://open.e.189.cn/",
	}

	data := map[string]string{
		"appKey":       "cloud",
		"accountType":  "01",
		"userName":     "{RSA}" + rsaEncrypt(username, key),
		"password":     "{RSA}" + rsaEncrypt(password, key),
		"validateCode": "",
		"captchaToken": captchaToken,
		"returnUrl":    returnUrl,
		"mailSuffix":   "@189.cn",
		"paramId":      paramId,
		"dynamicCheck": "FALSE",
		"clientType":   "10010",
		"cb_SaveName":  "1",
		"isOauth2":     "false",
	}

	resp, err = client.R().SetHeaders(headers).SetFormData(data).Post(url)
	if err != nil {
		log.Println(err)
		return false
	}

	msg, err := jsonparser.GetString(resp.Body(), "msg")
	if err != nil && msg != "登录成功" {
		log.Println(err, msg)
		return false
	}
	url, err = jsonparser.GetString(resp.Body(), "toUrl")
	if err != nil {
		log.Println(err)
		return false
	}
	client.R().Get(url)
	return true
}

// Do 执行
func Do(username, password string) string {
	client := resty.New()
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Linux; Android 5.1.1; SM-G930K Build/NRD90M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/74.0.3729.136 Mobile Safari/537.36 Ecloud/8.6.3 Android/22 clientId/355325117317828 clientModel/SM-G930K imsi/460071114317824 clientChannelId/qq proVersion/1.0.6",
	}
	client.SetHeaders(headers)

	if !login(client, username, password) {
		return "登录失败"
	}

	// 签到
	netdiskBonus := sign(client)

	return fmt.Sprintf("本次签到获得%dM空间", netdiskBonus)
}

func regexpString(reg, ctx string) string {
	r := regexp.MustCompile(reg)
	out := r.FindStringSubmatch(ctx)
	if len(out) != 2 {
		return ""
	}
	return out[1]
}

// rsaEncrypt 公钥加密
func rsaEncrypt(in, key string) string {
	//解密pem格式的公钥
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return ""
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Println(err)
		return ""
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(in))
	if err != nil {
		log.Println(err)
		return ""
	}
	return hex.EncodeToString(ciphertext)
}
