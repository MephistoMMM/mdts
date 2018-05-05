package handlers

import (
	"context"
	"fmt"
	"log"
	s2t "mdts/brokerSDK/s2t/broker"
	t2s "mdts/brokerSDK/t2s/broker"
	"mdts/dts/request"
	bmsg "mdts/protocols/brokermsg"
	pts "mdts/protocols/dtsproto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

const (
	s2tBrokerAddr = "127.0.0.1:9100"
	t2sBrokerAddr = "127.0.0.1:9110"
)

// TransforDataToThird ...
func TransforDataToThird(TID string, APICODE string, Data []byte) *bmsg.ResultToThird {

	conn, err := grpc.Dial(s2tBrokerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := s2t.NewS2TBrokerClient(conn)
	result, err := client.TransforDataToThird(context.Background(), &bmsg.ParamToThird{
		TID:     TID,
		APICODE: APICODE,
		Data:    Data,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return result
}

// TransforDataFromThird ...
func TransforDataFromThird(TID string, APICODE string, Data []byte) *bmsg.ResultFromThird {

	conn, err := grpc.Dial(s2tBrokerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := s2t.NewS2TBrokerClient(conn)
	result, err := client.TransforDataFromThird(context.Background(), &bmsg.ParamFromThird{
		TID:     TID,
		APICODE: APICODE,
		Data:    Data,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return result
}

// TransforDataToService ...
func TransforDataToService(Version string, APICODE string, Data []byte) *bmsg.ResultToService {

	conn, err := grpc.Dial(t2sBrokerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := t2s.NewT2SBrokerClient(conn)
	result, err := client.TransforDataToService(context.Background(), &bmsg.ParamToService{
		Version: Version,
		APICODE: APICODE,
		Data:    Data,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return result
}

// TransforDataFromService ...
func TransforDataFromService(Version string, APICODE string, Data []byte) *bmsg.ResultFromService {

	conn, err := grpc.Dial(t2sBrokerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := t2s.NewT2SBrokerClient(conn)
	result, err := client.TransforDataFromService(context.Background(), &bmsg.ParamFromService{
		Version: Version,
		APICODE: APICODE,
		Data:    Data,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return result
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

	result := TransforDataToThird(head.TID, head.APICODE, body)
	if result.GetState() == bmsg.EnumState_FAILED {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Broker(1) Error!",
		})
		return
	}

	_, byt, err := request.PostBytes(result.GetURL(), result.GetBody())
	if err != nil {
		errStr := fmt.Sprintf("Invalid Response From Third: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	result2 := TransforDataFromThird(head.TID, head.APICODE, byt)
	if result2.GetState() == bmsg.EnumState_FAILED {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Broker(2) Error!",
		})
		return
	}

	c.Data(200, "application/json", result2.GetBody())
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

	result := TransforDataToService(head.Version, head.APICODE, body)
	if result.GetState() == bmsg.EnumState_FAILED {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Broker(1) Error!",
		})
		return
	}

	_, byt, err := request.PostBytes(result.GetURL(), result.GetBody())
	if err != nil {
		errStr := fmt.Sprintf("Invalid Response From Service: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	result2 := TransforDataFromService(head.Version, head.APICODE, byt)
	if result2.GetState() == bmsg.EnumState_FAILED {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Broker(2) Error!",
		})
		return
	}

	c.Data(200, "application/json", result2.GetBody())
}
