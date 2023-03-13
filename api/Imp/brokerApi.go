package apiImp

import (
	"context"
	"log"
	"net"
	"therealbroker/api/proto"
	internalBroker "therealbroker/internal/broker"
	"therealbroker/pkg/broker"
	"therealbroker/pkg/prometheus"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedBrokerServer
	broker broker.Broker
}

func (s *server) Publish(ctx context.Context, in *proto.PublishRequest) (*proto.PublishResponse, error) {
	prometheus.MethodCalls.WithLabelValues("Publish").Inc()
	currentTime := time.Now()
	defer prometheus.MethodDuration.WithLabelValues("Publish").Observe(float64(time.Since(currentTime).Nanoseconds()))
	result, err := s.broker.Publish(ctx, in.Subject, broker.Message{
		Expiration: time.Second * time.Duration(in.ExpirationSeconds), Body: string(in.Body),
	})
	if err != nil {
		prometheus.MethodError.WithLabelValues("Publish").Inc()
		log.Println(err)
	}
	return &proto.PublishResponse{Id: int32(result)}, err
}

func (s *server) Subscribe(in *proto.SubscribeRequest, srv proto.Broker_SubscribeServer) error {
	prometheus.MethodCalls.WithLabelValues("Subscribe").Inc()
	prometheus.ActiveSubscribers.Inc()
	defer prometheus.ActiveSubscribers.Dec()
	ch, err := s.broker.Subscribe(nil, in.Subject)

	if err != nil {
		prometheus.MethodError.WithLabelValues("Subscribe").Inc()
		log.Println(err)
		return err
	}

	for message := range ch {
		msg := &proto.MessageResponse{Body: []byte(message.Body)}
		err := srv.Send(msg)
		if err != nil {
			prometheus.MethodError.WithLabelValues("Subscribe").Inc()
			return err
		}
	}
	return nil
}

func (s *server) Fetch(ctx context.Context, in *proto.FetchRequest) (*proto.MessageResponse, error) {
	prometheus.MethodCalls.WithLabelValues("Fetch").Inc()
	currentTime := time.Now()
	defer prometheus.MethodDuration.WithLabelValues("Fetch").Observe(float64(time.Since(currentTime).Nanoseconds()))
	result, err := s.broker.Fetch(ctx, in.Subject, int(in.Id))
	if err != nil {
		prometheus.MethodError.WithLabelValues("Fetch").Inc()
		log.Println(err)
		return nil, err
	}
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
