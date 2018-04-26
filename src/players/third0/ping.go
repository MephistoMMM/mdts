// Package main provide ping third
//
// Author: Mephis Pheies <mephistommm@gmail.com>
package main

import (
	"encoding/json"
	"log"
	"os"
	pts "players/protocols/dtsproto"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

const (
	dtsAddress = "http://127.0.0.1:7081/v1/"
	hostport   = ":9000"
	tid        = "a3x77n02"
)

var req = gorequest.New()

func receiveCancelOrder(c *gin.Context) {
	var body struct {
		OrderCode string `json:"orderCode"`
		Remark    string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		log.Fatalf("Receive Cancel Order Error: %v.", err)
	}

	log.Printf("Receive Cancel Order: {code: '%s', remark: '%s'}.",
		body.OrderCode, body.Remark)

	go func() {

		time.Sleep(time.Millisecond * 500)

		data := struct {
			OrderCode  string   `json:"orderCode"`
			ReasonType []string `json:"reasonType"`
			Remark     string   `json:"remark"`
		}{
			OrderCode:  body.OrderCode,
			ReasonType: []string{"999", "666"},
			Remark:     "Just A Test",
		}

		_, body, errs := req.Post(dtsAddress+"t2s").Type("json").
			Set("Version", "1.0.0").
			Set("APICODE", "00000003").
			Set("Sender", tid).
			SendStruct(&data).EndBytes()
		if errs != nil {
			log.Fatalf("Ping Error: %v.", errs)
		}

		var resData pts.CommResp
		if err := json.Unmarshal(body, &resData); err != nil {
			log.Fatalf("Send Refuse Order Error: %v.\n", err)
		}

		if resData.Code != pts.SUCCESS {
			log.Fatalf("Send Refuse Order Error: FAILED, %s.\n", resData.Message)
		}

		log.Printf("Refuse Order Success: %s .\n", string(resData.Data))
		os.Exit(0)
	}()

	c.JSON(200, &pts.CommResp{
		Code:    pts.SUCCESS,
		Message: "",
		Data:    []byte("{}"),
	})
}

func main() {
	router := gin.Default()
	router.POST("/cancelOrder", receiveCancelOrder)

	if err := router.Run(hostport); err != nil {
		log.Fatalln("Failed to run server: ", err)
	}
}
