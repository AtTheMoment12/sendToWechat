package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

// GetRequest Manage the HTTP GET request parameters
type GetRequest struct {
	urls url.Values
}

// Init Initializer
func (p *GetRequest) Init() *GetRequest {
	p.urls = url.Values{}
	return p
}

// InitFrom Initialized from another instance
func (p *GetRequest) InitFrom(reqParams *GetRequest) *GetRequest {
	if reqParams != nil {
		p.urls = reqParams.urls
	} else {
		p.urls = url.Values{}
	}
	return p
}

// AddParam Add URL escape property and value pair
func (p *GetRequest) AddParam(property string, value string) *GetRequest {
	if property != "" && value != "" {
		p.urls.Add(property, value)
	}
	return p
}

// BuildParams Concat the property and value pair
func (p *GetRequest) BuildParams() string {
	return p.urls.Encode()
}

func logInit() {
	f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("[info]:")
	if err != nil {
		log.Println(err)
	}
}
func getAccessToken() interface{} {
	url := "http://api.weixin.qq.com/cgi-bin/token?"
	init := new(GetRequest).Init()
	params := init.AddParam("grant_type", "client_credential").AddParam("appid", AppID).AddParam("secret", AppSecret).BuildParams()
	url = url + params
	req, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer req.Body.Close()
	dataMap := make(map[string]interface{})
	b, err := ioutil.ReadAll(req.Body)
	err = json.Unmarshal(b, &dataMap)
	if err != nil {
		log.Println(err)
	}
	return dataMap["access_token"]
}

func InitWxMpTemplateData() map[string]interface{} {
	data := make(map[string]interface{})
	return data
}
func addData(data map[string]interface{}, name string, value string, color string) {
	valueData := make(map[string]interface{})
	valueData["value"] = value
	valueData["color"] = color
	data[name] = valueData
}

func wechatPush(accessToken string, openId string, templateId string, data map[string]interface{}) {
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?"
	init := new(GetRequest).Init()
	params := init.AddParam("access_token", accessToken).BuildParams()
	url = url + params
	body := make(map[string]interface{})
	body["touser"] = openId
	body["template_id"] = templateId
	body["topcolor"] = "#FF0000"
	body["msgType"] = "event"
	body["data"] = data
	log.Println(body)

	configData, _ := json.Marshal(body)
	param := bytes.NewBuffer([]byte(configData))
	//fmt.Println(param)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, param)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println(err)
	}
	resp, err := client.Do(req)
	defer req.Body.Close()

	dataMap := make(map[string]interface{})
	b, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(b, &dataMap)
	if err != nil {
		log.Println(err)
	}

	if dataMap["errcode"] == 0.0 {
		if dataMap["errmsg"] == "ok" {
			log.Println("推送成功")
		} else {
			log.Println("推送失败，状态码正确", dataMap)
		}
	} else {
		log.Println("推送失败", dataMap)
	}
}

func getWeather() map[string]interface{} {
	url := "https://api.map.baidu.com/weather/v1/?district_id=410102&data_type=all&ak=CQLn4Mqi7fRikEAaXxeVaN0OvVdbWkmt"
	req, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer req.Body.Close()
	dataMap := make(map[string]interface{})
	b, err := ioutil.ReadAll(req.Body)
	err = json.Unmarshal(b, &dataMap)
	if err != nil {
		log.Println(err)
	}

	if dataMap["status"] == 0.0 {
		if dataMap["message"] == "success" {
			log.Println("返回成功")
			return dataMap["result"].(map[string]interface{})
		} else {
			log.Println("状态码正确，返回失败")
			return nil
		}
	} else {
		log.Println("返回失败")
		return nil
	}

}
