// Package main provide ping third
//
// Author: Mephis Pheies <mephistommm@gmail.com>
package main

import (
	"encoding/json"
	"log"
	pts "mdts/protocols/req2dts"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

const (
	dtsAddress   = "http://127.0.0.1:8081/v1/t2s"
	pingHostport = ":9000"
)

var (
	req           = gorequest.New()
	pingPongCount = 5
	hasBall       = true
)

/* Ping-Pong protocol
*
* http post /thirdping
* Content-Type: application/json
*
* Request:
*   | Key  | Type   |
*   | -----+------- |
*   | Ping | Number |
*
* Response:
*   | Key  | Type   |
*   | -----+------- |
*   | RePing | Number |
*
* Request:
*   | Key  | Type   |
*   | -----+------- |
*   | Pong | Number |
*
* Response:
*   | Key  | Type   |
*   | -----+------- |
*   | RePong | Number |
 */

type ppiReq struct {
	Ping int `json:"Ping"`
}

type ppiResp struct {
	RePing int `json:"RePing"`
}

type ppoReq struct {
	Pong int `json:"Pong"`
}

type ppoResp struct {
	RePong int `json:"RePong"`
}

func pong(c *gin.Context) {
	// Pong
	var reqData pts.RespS2T
	if err := c.ShouldBindJSON(&reqData); err != nil {
		log.Fatalf("Pong Error: %v.", err)
	}

	if reqData.Code != pts.SUCCESS {
		log.Fatalf("Pong Error: FAILED, %s.", reqData.Message)
	}

	var reqppdata ppoReq
	if err := json.Unmarshal([]byte(reqData.Data), &reqppdata); err != nil {
		log.Fatalf("Pong Error: %v.", err)
	}

	// Show Pong Data
	log.Printf("Pong: %d.", reqppdata.Pong)
	hasBall = true

	go func() {

		time.Sleep(time.Millisecond * 500)

		// check to end
		if pingPongCount < 0 {
			os.Exit(0)
		}

		// Ping
		ppdata := ppiReq{
			Ping: pingPongCount,
		}

		byt, _ := json.Marshal(&ppdata)
		data := pts.BodyT2S{
			Data: string(byt),
		}

		_, _, errs := req.Post(dtsAddress).Type("json").
			Set("TID", "a3x77n02UI3YWnhqr45UBe4AMCCq65NN").
			SendStruct(data).EndBytes()
		if errs != nil {
			log.Fatalf("Ping Error: %v.", errs)
		}
		hasBall = false
		pingPongCount--
	}()

	respppdata := ppoResp{
		RePong: reqppdata.Pong,
	}

	c.JSON(200, &respppdata)
}

func main() {
	router := gin.Default()
	router.POST("/thirdpong", pong)

	go func() {
		// Ping
		ppdata := ppiReq{
			Ping: pingPongCount,
		}

		byt, _ := json.Marshal(&ppdata)
		data := pts.BodyT2S{
			Data: string(byt),
		}

		_, body, errs := req.Post(dtsAddress).Type("json").
			Set("TID", "a3x77n02UI3YWnhqr45UBe4AMCCq65NN").
			SendStruct(data).EndBytes()
		if errs != nil {
			log.Fatalf("Ping Error: %v.", errs)
		}
		hasBall = false
		pingPongCount--

		var reqData pts.RespS2T
		if err := json.Unmarshal(body, &reqData); err != nil {
			log.Fatalf("Pong Error: %v.", err)
		}

		if reqData.Code != pts.SUCCESS {
			log.Fatalf("Pong Error: FAILED, %s.", reqData.Message)
		}

		log.Println(reqData.Data)
	}()

	if err := router.Run(pingHostport); err != nil {
		log.Fatalln("Failed to run server: ", err)
	}
}
