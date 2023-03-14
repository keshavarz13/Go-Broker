# Go Broker: A childish broker implemented by golang!

This message broker is a software that enables programs, systems and services to communicate with each other and exchange information over gRPC. 
Broker can have several topics and each message published to certain topic will be broadcasted
to all subscribers to that topic.

Features of this broker:
- The possibility of saving messages using the `PostgreSQL`
- The possibility of using inMemory mode
- Ability to transmit more than 21 thousand messages per second in persistent mode
- Using `Prometheus` to measure broker performance metrics

The broker run on port 5100 and prometheus metrics can be accessed from :9000/metrics

## RPCs Description
- Publish Requst
```protobuf
message PublishRequest {
  string subject = 1;
  bytes body = 2;
  int32 expirationSeconds = 3;
}
```
- Fetch Request
```protobuf
message FetchRequest {
  string subject = 1;
  int32 id = 2;
}
```
- Subscribe Request
```protobuf
message SubscribeRequest {
  string subject = 1;
}
```
- RPC Service
```protobuf
service Broker {
  rpc Publish (PublishRequest) returns (PublishResponse);
  rpc Subscribe(SubscribeRequest) returns (stream MessageResponse);
  rpc Fetch(FetchRequest) returns (MessageResponse);
}
```

# How to Run it?
You can directly use Docker to run this broker:
```shell
docker build -t my-broker . 
docker run -p 9000:9000 -p 5100:5100 my-broker
```
Or you can run the broker by running the following command:
```shell
chmod +x run.sh
./run.sh
```
