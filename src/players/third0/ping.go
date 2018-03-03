// Package main provide ping third
//
// Author: Mephis Pheies <mephistommm@gmail.com>
package main

import (
	"encoding/json"
	"log"
	"os"
	pts "players/protocols/pingpong"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

const (
	dtsAddress = "http://127.0.0.1:7081/v1/"
	hostport   = ":9000"
	tid        = "a3x77n02UI3YWnhqr45UBe4AMCCq65NN"
)

var (
	req           = gorequest.New()
	pingPongCount = 5
	hasBall       = true
)

func receivePong(c *gin.Context) {
	// Pong
	var body pts.BodyPongS2T

	if err := c.ShouldBindJSON(&body); err != nil {
		log.Fatalf("Ping Error: %v.", err)
	}

	// Show Pong Data
	log.Printf("Receive Pong Ball(%s).", body.Ball)
	hasBall = true

	go func() {

		time.Sleep(time.Millisecond * 500)

		// check to end
		if pingPongCount < 0 {
			os.Exit(0)
		}

		// Ping
		ball := strconv.Itoa(pingPongCount)
		data := pts.BodyPingT2S{
			Ball:   ball,
			Time:   time.Now().UnixNano(),
			Sender: tid,
		}

		_, body, errs := req.Post(dtsAddress+"pingt2s").Type("json").
			Set("Version", "1.0.0").
			SendStruct(data).EndBytes()
		if errs != nil {
			log.Fatalf("Ping Error: %v.", errs)
		}
		hasBall = false
		pingPongCount--

		var reqData pts.CommResp
		if err := json.Unmarshal(body, &reqData); err != nil {
			log.Fatalf("Response Ping Error: %v.", err)
		}

		if reqData.Code != pts.SUCCESS {
			log.Fatalf("Response Ping Error: FAILED, %s.", reqData.Message)
		}

		log.Printf("Ping Ball(%s) Success.", reqData.Data)
	}()

	c.JSON(200, &pts.CommResp{
		Code:    pts.SUCCESS,
		Message: "",
		Data:    body.Ball,
	})
}

func main() {
	router := gin.Default()
	router.POST("/pongs2t", receivePong)

	go func() {
		// Ping
		ball := strconv.Itoa(pingPongCount)
		data := pts.BodyPingT2S{
			Ball:   ball,
			Time:   time.Now().UnixNano(),
			Sender: tid,
		}

		_, body, errs := req.Post(dtsAddress+"pingt2s").Type("json").
			Set("Version", "1.0.0").
			SendStruct(data).EndBytes()
		if errs != nil {
			log.Fatalf("Ping Error: %v.", errs)
		}
		hasBall = false
		pingPongCount--

		var reqData pts.CommResp
		if err := json.Unmarshal(body, &reqData); err != nil {
			log.Fatalf("Response Ping Error: %v.", err)
		}

		if reqData.Code != pts.SUCCESS {
			log.Fatalf("Response Ping Error: FAILED, %s.", reqData.Message)
		}

		log.Printf("Ping Ball(%s) Success.", reqData.Data)
	}()

	if err := router.Run(hostport); err != nil {
		log.Fatalln("Failed to run server: ", err)
	}
}
