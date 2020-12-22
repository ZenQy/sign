package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/wxpusher/wxpusher-sdk-go"
	"github.com/wxpusher/wxpusher-sdk-go/model"
	"github.com/zenqy/sign/mt"
	"github.com/zenqy/sign/paoluz"
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
			switch k {
			case "paoluz":
				txt := paoluz.Sign(conf[k][i])
				msg += fmt.Sprintf("%s | %s | %s\n", k, conf[k][i]["username"], txt)
			case "mt":
				txt := mt.Sign(conf[k][i])
				msg += fmt.Sprintf("%s | %s | %s\n", k, conf[k][i]["username"], txt)
			}

		}
	}

	data, err = yaml.Marshal(&conf)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(fn, data, 0644)
	sendMsg(msg)
}

func sendMsg(content string) error {
	appToken := "AT_JSR0rKIq46EzzP9nt7wjY6p5Htc4JJy3"
	uid := "UID_JeT9csRRPGrm8xfLWuwJH0CZelYd"
	msg := model.NewMessage(appToken).SetSummary("签到情况：" + time.Now().Format("2006-01-02")).SetContentType(3).SetContent(content).AddUId(uid)
	_, err := wxpusher.SendMessage(msg)
	return err
}
