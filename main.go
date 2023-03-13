package main

import (
	"context"
	"log"
	"net"
	"therealbroker/api/proto"
	"therealbroker/internal/broker"
	pkgBroker "therealbroker/pkg/broker"
	"time"

	"google.golang.org/grpc"
)

// Main requirements:
// 1. All tests should be passed
// 2. Your logs should be accessible in Graylog
// 3. Basic prometheus metrics ( latency, throughput, etc. ) should be implemented
// 	  for every base functionality ( publish, subscribe etc. )

type server struct {
	proto.UnimplementedBrokerServer
	broker pkgBroker.Broker
}

func (s *server) Publish(ctx context.Context, in *proto.PublishRequest) (*proto.PublishResponse, error) {
	result, err := s.broker.Publish(ctx, in.Subject, pkgBroker.Message{
		Expiration: time.Second * time.Duration(in.ExpirationSeconds), Body: string(in.Body),
	})
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
	return nil
}

func (s *server) Fetch(ctx context.Context, in *proto.FetchRequest) (*proto.MessageResponse, error) {
	result, err := s.broker.Fetch(ctx, in.Subject, int(in.Id))
	log.Println(err)
	return &proto.MessageResponse{Body: []byte(result.Body)}, err
}

func main() {
	listener, err := net.Listen("tcp", ":5100")
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	s := grpc.NewServer()
	myServer := server{broker: broker.NewModule()}
	proto.RegisterBrokerServer(s, &myServer)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
