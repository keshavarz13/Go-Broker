package main

import (
	"context"
	"log"
	api "therealbroker/api/proto"
	"time"

	"google.golang.org/grpc"
)

const VUs = 1000
const REQUESTS = 1000

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func Publish(client api.BrokerClient) {
	_, err := client.Publish(context.Background(), &api.PublishRequest{
		Subject:           "zzzzz",
		Body:              []byte("hi my dear friends!"),
		ExpirationSeconds: 2000,
	})
	if err != nil {
		log.Println("Error publishing message: ", err)
		return
	}
}

func Fetch(client api.BrokerClient) {
	_, err := client.Fetch(context.Background(), &api.FetchRequest{
		Subject: "zzzzz",
		Id:      2,
	})
	if err != nil {
		log.Println("Error fetching message: ", err)
		return
	}
}

func main() {

	conn, err := grpc.Dial("localhost:5100", grpc.WithInsecure())
	if err != nil {
		log.Println("Error connecting to broker: ", err)
		return
	}
	defer conn.Close()

	client := api.NewBrokerClient(conn)

	for i := 0; i < VUs; i++ {
		go func() {
			for j := 0; j < REQUESTS; j++ {
				Publish(client)
				// Fetch(client)
			}
		}()
	}

	<-time.After(time.Minute * 10)
}
