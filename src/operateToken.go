package main

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strings"
)

func echoServer(router *gin.Engine) {
	router.POST("/api/checkToken", func(c *gin.Context) {
		//fmt.Println(c.PostForm("errcode"), c.Query("errmg"), c.Query("msgid"))
		c.String(http.StatusOK, "success")
	})
}

//验证令牌
func checkToken(router *gin.Engine) {
	router.GET("/api/checkToken", func(c *gin.Context) {
		signature := c.Query("signature")
		echostr := c.Query("echostr")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		//根据微信文档要求加入数组并排序
		var arr []string
		arr = append(arr, timestamp, nonce, WechatToken)
		sort.Strings(arr)
		str := strings.Join(arr, ``)
		//sha1加密
		s := sha1.New()
		s.Write([]byte(str))
		bs := s.Sum(nil)
		compare := hex.EncodeToString(bs)
		//比对令牌
		if compare == signature {
			c.String(http.StatusOK, echostr)
		} else {
			c.String(http.StatusBadRequest, "error,checksum failed")
		}
	})
}
