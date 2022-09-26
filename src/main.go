package main

import (
	"github.com/robfig/cron"
	"log"
	"strconv"
	"time"
)

var WechatToken = "411a9133e3e577588aaf38129ab748d866cf8840"
var AppID = "wxd392c9a557409b79"
var AppSecret = "e4f47224c831a70e10d5c412ef62cdca"
var templateOfGoodMorning = "ppBXCx3_eWp_5xQeoXZsadLJUzbvU6OZt7UzlriXTGY"
var testOpenID = "o_nML68aPWeGx3HnR3e_V_4ZWJtM"
var openID = "o_nML6-QCyAKkPpDFjzw0QxSK3oA"

func test() {

}

func main() {
	//router := gin.New()
	//checkToken(router)
	//echoServer(router)
	//router.Run(":80")
	logInit()
	log.Println("service start")
	c := cron.New()
	c.AddFunc("0 30 9 * * ?", func() {
		log.Println("start push ")
		pushEveryday(openID)
		//pushEveryday(testOpenID)
		log.Println("push success")
	})
	c.Start()
	select {}
}

func pushEveryday(openId string) {
	accessToken := getAccessToken()

	weatherMap := getWeather()
	log.Println(weatherMap)
	weather := weatherMap["now"].(map[string]interface{})["text"].(string)
	low := strconv.FormatFloat(weatherMap["forecasts"].([]interface{})[0].(map[string]interface{})["low"].(float64), 'f', 0, 64)
	high := strconv.FormatFloat(weatherMap["forecasts"].([]interface{})[0].(map[string]interface{})["high"].(float64), 'f', 0, 64)

	timeNow, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	startTime, _ := time.Parse("2006-01-02", "2015-08-31")
	dis := timeNow.Sub(startTime)
	love := dis.Hours()/24 + 1
	var remark string
	if time.Now().Format("01-02") == "08-31" {
		//current, _ := strconv.Atoi(time.Now().Format("2006"))
		//year := 2015
		remark = "今天是恋爱纪念日~"
	} else if time.Now().Format("01-02") == "02-14" {
		remark = "今天是情人节~"
	} else {
		remark = "新的一天要快乐哦~"
	}
	data := InitWxMpTemplateData()
	addData(data, "date", time.Now().Format("2006-01-02")+" "+time.Now().Weekday().String(), "#173177")
	addData(data, "remark", remark, "#173177")
	addData(data, "weather", weather, "#173177")
	addData(data, "low", low, "#173177")
	addData(data, "high", high, "#173177")
	addData(data, "love", strconv.FormatFloat(love, 'f', 0, 64), "#173177")
	wechatPush(accessToken.(string), openId, templateOfGoodMorning, data)
}
