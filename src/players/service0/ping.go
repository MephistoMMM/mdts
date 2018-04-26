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
	dtsAddress = "http://127.0.0.1:7080/v1/"
	hostport   = ":9100"
	version    = "1.0.0"
)

var req = gorequest.New()

func sendCancelOrder() {
	data := `{"orderCode":"12345678", "remark":"Just A Test"}`

	_, body, errs := req.Post(dtsAddress+"s2t").Type("json").
		Set("Sender", "1.0.0").
		Set("TID", "a3x77n02").
		Set("APICODE", "00000002").
		Send(data).EndBytes()
	if errs != nil {
		log.Fatalf("Ping Error: %v.", errs)
	}

	var resData pts.CommResp
	if err := json.Unmarshal(body, &resData); err != nil {
		log.Fatalf("Send Cancel Order Error: %v. \n", err)
	}
	if resData.Code != pts.SUCCESS {
		log.Fatalf("Send Cancel Order Error: FAILED, %s.\n", resData.Message)
	}

	log.Printf("Cancel Order Success: %s .\n", string(resData.Data))
}

func receiveRefuseOrder(c *gin.Context) {
	var body struct {
		OrderCode  string   `json:"orderCode"`
		ReasonType []string `json:"reasonType"`
		Remark     string   `json:"remark"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		log.Fatalf("Receive Refuse Order Error: %v.", err)
	}

	// Show Ping Data
	log.Printf("Receive Refuse Order: {code: '%s', reasonType: %v, remark: '%s'}.",
		body.OrderCode, body.ReasonType, body.Remark)

	go func() {
		time.Sleep(2000 * time.Millisecond)
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
	router.POST("/rescue/refuseOrder", receiveRefuseOrder)

	go func() {
		time.Sleep(500 * time.Millisecond)
		sendCancelOrder()
	}()

	if err := router.Run(hostport); err != nil {
		log.Fatalln("Failed to run server: ", err)
	}
}
