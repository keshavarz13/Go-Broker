package apiImp

import (
	"context"
	"log"
	"net"
	"therealbroker/api/proto"
	internalBroker "therealbroker/internal/broker"
	"therealbroker/pkg/broker"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedBrokerServer
	broker broker.Broker
}

func (s *server) Publish(ctx context.Context, in *proto.PublishRequest) (*proto.PublishResponse, error) {
	result, err := s.broker.Publish(ctx, in.Subject, broker.Message{
		Expiration: time.Second * time.Duration(in.ExpirationSeconds), Body: string(in.Body),
	})
	log.Println(err)
	return &proto.PublishResponse{Id: int32(result)}, err
}

func (s *server) Subscribe(in *proto.SubscribeRequest, srv proto.Broker_SubscribeServer) error {
	ch, _ := s.broker.Subscribe(nil, in.Subject)
	for {
		select {
		case msg := <-ch:
			srv.Send(&proto.MessageResponse{Body: []byte(msg.Body)})
		}
	}
}

func (s *server) Fetch(ctx context.Context, in *proto.FetchRequest) (*proto.MessageResponse, error) {
	result, err := s.broker.Fetch(ctx, in.Subject, int(in.Id))
	log.Println(err)
	return &proto.MessageResponse{Body: []byte(result.Body)}, err
}

func RunGrpcServer() {
	listener, err := net.Listen("tcp", ":5100")
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	s := grpc.NewServer()
	myServer := server{broker: internalBroker.NewModule()}
	proto.RegisterBrokerServer(s, &myServer)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
