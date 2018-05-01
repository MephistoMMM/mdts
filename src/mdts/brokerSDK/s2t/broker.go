package s2t

import (
	"context"
	"errors"
	"log"
	pb "mdts/brokerSDK/broker"
	bmsg "mdts/protocols/brokermsg"
	"net"

	"google.golang.org/grpc"
)

var (
	errConflictOfID = errors.New("Error Conflict Of TID")
)

// Server provide a method to run a grpc server
type Server struct {
	t Transformer
}

// TransforDataToThird implement service S2TBroker.TransforDataToThird
func (svr *Server) TransforDataToThird(ctx context.Context, p *bmsg.ParamToThird) (*bmsg.ResultToThird, error) {
	if svr.t.ID() != p.GetTID() {
		return nil, errConflictOfID
	}

	tfr, err := svr.t.TransTo(p.GetAPICODE(), p.GetData())
	if err != nil {
		return nil, err
	}

	return &bmsg.ResultToThird{
		State:  bmsg.EnumState_SUCCESS,
		Method: tfr.Method,
		Head:   tfr.Head,
		Body:   tfr.Body,
		URL:    tfr.URL,
	}, nil
}

// TransforDataFromThird implement service S2TBroker.TransforDataFromThird
func (svr *Server) TransforDataFromThird(ctx context.Context, p *bmsg.ParamFromThird) (*bmsg.ResultFromThird, error) {
	if svr.t.ID() != p.GetTID() {
		return nil, errConflictOfID
	}

	tfr, err := svr.t.TransFrom(p.GetAPICODE(), p.GetData())
	if err != nil {
		return nil, err
	}

	return &bmsg.ResultFromThird{
		State: bmsg.EnumState_SUCCESS,
		Head:  tfr.Head,
		Body:  tfr.Body,
	}, nil
}

// Run run server
func (svr *Server) Run(hostport string, t Transformer) error {
	svr.run(hostport, t)
}

func (svr *Server) run(hostport string, t Transformer) {
	svr.t = t

	lis, err := net.Listen("tcp", hostport)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterS2TBrokerServer(grpcServer, svr)
	grpcServer.Serve(lis)
}
