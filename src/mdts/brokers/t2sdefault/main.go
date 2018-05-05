package main

import (
	"log"
	bsdk "mdts/brokerSDK/base"
	"mdts/brokerSDK/t2s"
	bmsg "mdts/protocols/brokermsg"
)

type defaultTrans struct {
	id  string
	url string
}

func NewDefaultTrans(id string, url string) *defaultTrans {
	return &defaultTrans{
		id:  id,
		url: url,
	}
}

func (dt *defaultTrans) ID() string {
	return dt.id
}

func (dt *defaultTrans) TransTo(APICODE string, Data []byte) (*bsdk.TransToResult, error) {
	log.Println(string(Data))
	return &bsdk.TransToResult{
		Method: bmsg.EnumMethod_HttpPost,
		Head:   make(map[string]string),
		Body:   Data,
		URL:    dt.url,
	}, nil
}

func (dt *defaultTrans) TransFrom(APICODE string, Data []byte) (*bsdk.TransFromResult, error) {
	log.Println(string(Data))
	return &bsdk.TransFromResult{
		Head: make(map[string]string),
		Body: Data,
	}, nil
}

const (
	ID  = "1.0.0"
	URL = "http://127.0.0.1:9010/rescue/refuseOrder"

	hostport = ":9110"
)

func main() {
	trans := NewDefaultTrans(ID, URL)
	server := &t2s.Server{}

	if err := server.Run(hostport, trans); err != nil {
		log.Fatalln(err)
	}
}
