package t2s

import (
	"context"
	"errors"
	"log"
	"mdts/brokerSDK/base"
	pb "mdts/brokerSDK/t2s/broker"
	bmsg "mdts/protocols/brokermsg"
	"net"

	"google.golang.org/grpc"
)

var (
	errConflictOfID = errors.New("Error Conflict Of Version")
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

// TransforDataToService implement service T2SBroker.TransforDataToService
func (svr *Server) TransforDataToService(ctx context.Context, p *bmsg.ParamToService) (*bmsg.ResultToService, error) {
	if svr.t.ID() != p.GetVersion() {
		return nil, errConflictOfID
	}

	tfr, err := svr.t.TransTo(p.GetAPICODE(), p.GetData())
	if err != nil {
		return nil, err
	}

	return &bmsg.ResultToService{
		State:  bmsg.EnumState_SUCCESS,
		Method: tfr.Method,
		Head:   tfr.Head,
		Body:   tfr.Body,
		URL:    tfr.URL,
	}, nil
}

// TransforDataFromService implement service T2SBroker.TransforDataFromService
func (svr *Server) TransforDataFromService(ctx context.Context, p *bmsg.ParamFromService) (*bmsg.ResultFromService, error) {
	if svr.t.ID() != p.GetVersion() {
		return nil, errConflictOfID
	}

	tfr, err := svr.t.TransFrom(p.GetAPICODE(), p.GetData())
	if err != nil {
		return nil, err
	}

	return &bmsg.ResultFromService{
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
	pb.RegisterT2SBrokerServer(grpcServer, svr)
	return grpcServer.Serve(lis)
}
