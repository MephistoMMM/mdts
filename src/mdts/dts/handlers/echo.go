package handlers

import (
	"io/ioutil"
	pts "mdts/protocols/req2dts"

	"github.com/gin-gonic/gin"
)

// Echo ...
func Echo(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, "Failed read request body")
		return
	}
	c.JSON(200, &pts.RespT2S{
		Code:    pts.SUCCESS,
		Message: "",
		Data:    string(data),
	})
}

// PingPong ...
func PingPong(c *gin.Context) {

}
