// Package req2dts provide protocol of service request dts
//
// Author: Mephis Pheies <mephistommm@gmail.com>
package req2dts

import (
	pts "mdts/protocols"

	"github.com/gin-gonic/gin"
)

type respcode uint8

const (
	SUCCESS respcode = 0
	FAILED  respcode = 1
)

// CommResp is the common response struct
type CommResp struct {
	Code    respcode
	Message string
}

// HeadS2T is the Request Head from service to third
type HeadS2T struct {
	TID string `binding:"alphanum,len=32"`
}

// HeadS2T is the Request Body from service to third
type BodyS2T struct {
	Data string
}

// RespS2T is the response from service to third
type RespS2T struct {
	Code    respcode
	Message string
	Data    string
}

// HeadT2S is the Request Head from third to service
type HeadT2S struct {
	Version string `binding:"serversion"`
}

// HeadT2S is the Request Body from third to service
type BodyT2S struct {
	Data string
}

// RespT2S is the response from third to service
type RespT2S struct {
	Code    respcode
	Message string
	Data    string
}

// ParseReqS2T parse request from service to third
func ParseReqS2T(c *gin.Context) (*HeadS2T, *BodyS2T, error) {
	var (
		h HeadS2T
		b BodyS2T
	)
	h.TID = c.GetHeader("TID")

	if err := pts.Validator.ValidateStruct(h); err != nil {
		return nil, nil, err
	}

	if err := c.ShouldBindJSON(&b); err != nil {
		return nil, nil, err
	}

	return &h, &b, nil
}

// ParseReqS2T parse request from third to service
func ParseReqT2S(c *gin.Context) (*HeadT2S, *BodyT2S, error) {
	var (
		h HeadT2S
		b BodyT2S
	)
	h.Version = c.GetHeader("Version")

	if err := pts.Validator.ValidateStruct(h); err != nil {
		return nil, nil, err
	}

	if err := c.ShouldBindJSON(&b); err != nil {
		return nil, nil, err
	}

	return &h, &b, nil
}
