package handlers

import (
	"fmt"
	"log"
	"mdts/dts/request"
	pts "mdts/protocols/dtsproto"

	"github.com/gin-gonic/gin"
)

// TransforDataToThird ...
func TransforDataToThird(
	TID string,
	APICODE string,
	Data []byte) (state int, method int, head map[string]string, body []byte, url string) {

	return 0, 1, make(map[string]string), Data, thirdPath + "cancelOrder"
}

// TransforDataFromThird ...
func TransforDataFromThird(
	TID string,
	APICODE string,
	Data []byte) (state int, head map[string]string, body []byte) {

	return 0, make(map[string]string), Data
}

// TransforDataToService ...
func TransforDataToService(
	Version string,
	APICODE string,
	Data []byte) (state int, method int, head map[string]string, body []byte, url string) {

	return 0, 1, make(map[string]string), Data, servicePath + "rescue/refuseOrder"
}

// TransforDataFromService ...
func TransforDataFromService(
	Version string,
	APICODE string,
	Data []byte) (state int, head map[string]string, body []byte) {

	return 0, make(map[string]string), Data
}

// HandleS2T ...
func HandleS2T(c *gin.Context) {
	head, body, err := pts.ParseS2T(c)
	if err != nil {
		errStr := fmt.Sprintf("Invalid Request From Service: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	// Get Info Of TID From Etcd
	log.Printf("Request To TID: %s.", head.TID)

	state, _, _, bs, url := TransforDataToThird(head.TID, head.APICODE, body)
	if state == int(pts.ABEND) {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Broker(1) Error!",
		})
		return
	}

	_, byt, err := request.PostBytes(url, bs)
	if err != nil {
		errStr := fmt.Sprintf("Invalid Response From Third: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	state, _, byt = TransforDataFromThird(head.TID, head.APICODE, byt)
	if state == int(pts.ABEND) {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Broker(2) Error!",
		})
		return
	}

	c.Data(200, "application/json", byt)
}

// HandleT2S ...
func HandleT2S(c *gin.Context) {
	head, body, err := pts.ParseT2S(c)
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

	state, _, _, bs, url := TransforDataToService(head.Version, head.APICODE, body)
	if state == int(pts.ABEND) {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Broker(1) Error!",
		})
		return
	}

	_, byt, err := request.PostBytes(url, bs)
	if err != nil {
		errStr := fmt.Sprintf("Invalid Response From Service: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	state, _, byt = TransforDataFromService(head.Version, head.APICODE, byt)
	if state == int(pts.ABEND) {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Broker(2) Error!",
		})
		return
	}

	c.Data(200, "application/json", byt)
}
