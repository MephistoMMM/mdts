package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"mdts/dts/request"
	pts "mdts/protocols/pingpong"

	"github.com/gin-gonic/gin"
)

const (
	thirdPath   = "http://127.0.0.1:9000/"
	servicePath = "http://127.0.0.1:9100/"
)

// Echo ...
func Echo(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, "Failed read request body")
		return
	}
	c.JSON(200, &pts.CommResp{
		Code:    pts.SUCCESS,
		Message: "",
		Data:    string(data),
	})
}

// PingZeroForOut ...
func PingZeroForOut(c *gin.Context) {
	head, body, err := pts.ParsePingT2S(c)
	if err != nil {
		errStr := fmt.Sprintf("Invalid Request From Third: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	// Get Info Of Version From Etcd
	log.Printf("Request To Version: %s.", head.Version)

	_, byt, err := request.Post(servicePath+"pingt2s", &body)
	if err != nil {
		errStr := fmt.Sprintf("Invalid Response From Service: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	c.Data(200, "application/json", byt)
}

// PongZeroForOut ...
func PongZeroForOut(c *gin.Context) {
	head, body, err := pts.ParsePongT2S(c)
	if err != nil {
		errStr := fmt.Sprintf("Invalid Request From Third: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	// Get Info Of Version From Etcd
	log.Printf("Request To Version: %s.", head.Version)

	_, byt, err := request.Post(servicePath+"pongt2s", &body)
	if err != nil {
		errStr := fmt.Sprintf("Invalid Response From Service: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	c.Data(200, "application/json", byt)
}

// PingZeroForIn ..
func PingZeroForIn(c *gin.Context) {
	head, body, err := pts.ParsePingS2T(c)
	if err != nil {
		errStr := fmt.Sprintf("Invalid Request From Third: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	// Get Info Of TID From Etcd
	log.Printf("Request To TID: %s.", head.TID)

	_, byt, err := request.Post(thirdPath+"pings2t", &body)
	if err != nil {
		errStr := fmt.Sprintf("Invalid Response From Service: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	c.Data(200, "application/json", byt)
}

// PongZeroForIn ...
func PongZeroForIn(c *gin.Context) {
	head, body, err := pts.ParsePongS2T(c)
	if err != nil {
		errStr := fmt.Sprintf("Invalid Request From Third: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	// Get Info Of TID From Etcd
	log.Printf("Request To TID: %s.", head.TID)

	_, byt, err := request.Post(thirdPath+"pongs2t", &body)
	if err != nil {
		errStr := fmt.Sprintf("Invalid Response From Service: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	c.Data(200, "application/json", byt)
}
