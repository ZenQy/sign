package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/wxpusher/wxpusher-sdk-go"
	"github.com/wxpusher/wxpusher-sdk-go/model"
	"github.com/zenqy/sign/cloud189"
	"github.com/zenqy/sign/mt"
	"github.com/zenqy/sign/paoluz"
	"github.com/zenqy/sign/pdawiki"
	"gopkg.in/yaml.v2"
)

func main() {

	fn := "config.yaml"
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	conf := map[string]([]map[string]string){}
	if err := yaml.Unmarshal(data, &conf); err != nil {
		panic(err)
	}
	msg := "网站|帐号|签到\n:--:|:--:|:--:\n"
	for k := range conf {
		for i := range conf[k] {
			username, password := conf[k][i]["username"], conf[k][i]["password"]
			switch k {
			case "paoluz":
				txt := paoluz.Do(username, password)
				msg += fmt.Sprintf("%s | %s | %s\n", k, conf[k][i]["username"], txt)
			case "mt":
				txt := mt.Do(username, password)
				msg += fmt.Sprintf("%s | %s | %s\n", k, conf[k][i]["username"], txt)
			case "pdawiki":
				txt := pdawiki.Do(username, password)
				msg += fmt.Sprintf("%s | %s | %s\n", k, conf[k][i]["username"], txt)
			case "cloud189":
				txt := cloud189.Do(username, password)
				msg += fmt.Sprintf("%s | %s | %s\n", k, conf[k][i]["username"], txt)
			}

		}
	}

	sendMsg(msg)
	log.Println("本次运行完成...\n" + msg)
}

func sendMsg(content string) error {
	appToken := "AT_JSR0rKIq46EzzP9nt7wjY6p5Htc4JJy3"
	uid := "UID_JeT9csRRPGrm8xfLWuwJH0CZelYd"
	msg := model.NewMessage(appToken).SetSummary("签到情况：" + time.Now().Format("2006-01-02")).SetContentType(3).SetContent(content).AddUId(uid)
	_, err := wxpusher.SendMessage(msg)
	return err
}
