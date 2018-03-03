// Package main provide ping third
//
// Author: Mephis Pheies <mephistommm@gmail.com>
package main

import (
	"encoding/json"
	"log"
	"os"
	pts "players/protocols/pingpong"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

const (
	dtsAddress = "http://127.0.0.1:7080/v1/"
	hostport   = ":9100"
	version    = "1.0.0"
)

var (
	req           = gorequest.New()
	pingPongCount = 5
	hasBall       = false
)

func receivePing(c *gin.Context) {
	// Ping
	var body pts.BodyPingT2S

	if err := c.ShouldBindJSON(&body); err != nil {
		log.Fatalf("Ping Error: %v.", err)
	}

	// Show Ping Data
	log.Printf("Receive Ping Ball(%s).", body.Ball)
	hasBall = true

	go func() {

		time.Sleep(time.Millisecond * 500)

		// Pong
		data := pts.BodyPingT2S{
			Ball: body.Ball,
			Time: time.Now().UnixNano(),
		}

		_, body, errs := req.Post(dtsAddress+"pongs2t").Type("json").
			Set("TID", body.Sender).
			SendStruct(data).EndBytes()
		if errs != nil {
			log.Fatalf("Ping Error: %v.", errs)
		}
		hasBall = false

		var reqData pts.CommResp
		if err := json.Unmarshal(body, &reqData); err != nil {
			log.Fatalf("Response Ping Error: %v.", err)
		}

		if reqData.Code != pts.SUCCESS {
			log.Fatalf("Response Ping Error: FAILED, %s.", reqData.Message)
		}

		log.Printf("Pong Ball(%s) Success.", reqData.Data)
		pingPongCount--

		// check to end
		if pingPongCount < 0 {
			os.Exit(0)
		}
	}()

	c.JSON(200, &pts.CommResp{
		Code:    pts.SUCCESS,
		Message: "",
		Data:    body.Ball,
	})
}

func main() {
	router := gin.Default()
	router.POST("/pingt2s", receivePing)

	if err := router.Run(hostport); err != nil {
		log.Fatalln("Failed to run server: ", err)
	}
}
