package s2t

import (
	"context"
	"errors"
	"log"
	"mdts/brokerSDK/base"
	pb "mdts/brokerSDK/s2t/broker"
	bmsg "mdts/protocols/brokermsg"
	"net"

	"google.golang.org/grpc"
)

var (
	errConflictOfID = errors.New("Error Conflict Of TID")
	endpoints       = []string{"133.130.119.62:2379", "133.130.119.62:2381", "133.130.119.62:2383"}
)

// Server provide a method to run a grpc server
type Server struct {
	t base.Transformer
	s *base.Service
}

func NewServer(T string, ip string, port string) *Server {
	service, err := base.NewService(T, ip, port, endpoints)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	server := &Server{
		s: service,
	}

	go func() {
		log.Println("Discovery Service Start.")
		service.Start()
		log.Println("Discovery Service Stop.")
	}()

	return server
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
func (svr *Server) Run(hostport string, t base.Transformer) error {
	return svr.run(hostport, t)
}

func (svr *Server) run(hostport string, t base.Transformer) error {
	svr.t = t

	lis, err := net.Listen("tcp", hostport)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterS2TBrokerServer(grpcServer, svr)
	return grpcServer.Serve(lis)
}
