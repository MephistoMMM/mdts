// Package pingpong provide protocol of service request dts
//
// Author: Mephis Pheies <mephistommm@gmail.com>
package pingpong

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
	Data    string
}

// * PingS2T 00000001

// HeadPingS2T is the Request Head from service to third
type HeadPingS2T struct {
	TID string `binding:"alphanum,len=32"`
}

// HeadPingS2T is the Request Body from service to third
type BodyPingS2T struct {
	Ball   string
	Time   int64
	Sender string
}

// ParsePingS2T parse request from service to third
func ParsePingS2T(c *gin.Context) (*HeadPingS2T, *BodyPingS2T, error) {
	var (
		h HeadPingS2T
		b BodyPingS2T
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

// * PongS2T 00000002

// HeadPongS2T is the Request Head from service to third
type HeadPongS2T struct {
	TID string `binding:"alphanum,len=32"`
}

// HeadPongS2T is the Request Body from service to third
type BodyPongS2T struct {
	Ball string
	Time int64
}

// ParsePongS2T parse request from service to third
func ParsePongS2T(c *gin.Context) (*HeadPongS2T, *BodyPongS2T, error) {
	var (
		h HeadPongS2T
		b BodyPongS2T
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

// * PingT2S 10000001

// HeadPingT2S is the Request Head from third to service
type HeadPingT2S struct {
	Version string `binding:"serversion"`
}

// HeadPingT2S is the Request Body from third to service
type BodyPingT2S struct {
	Ball   string
	Time   int64
	Sender string
}

// ParsePingS2T parse request from third to service
func ParsePingT2S(c *gin.Context) (*HeadPingT2S, *BodyPingT2S, error) {
	var (
		h HeadPingT2S
		b BodyPingT2S
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

// * PongT2S 10000002

// HeadPongT2S is the Request Head from third to service
type HeadPongT2S struct {
	Version string `binding:"serversion"`
}

// HeadPongT2S is the Request Body from third to service
type BodyPongT2S struct {
	Ball string
	Time int64
}

// ParsePongS2T parse request from third to service
func ParsePongT2S(c *gin.Context) (*HeadPongT2S, *BodyPongT2S, error) {
	var (
		h HeadPongT2S
		b BodyPongT2S
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
