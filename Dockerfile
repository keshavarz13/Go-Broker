FROM golang:1.16-alpine

WORKDIR /app

COPY ./ ./
# RUN go mod download
RUN go build -o /my-broker

EXPOSE 9000
EXPOSE 5100

CMD [ "/my-broker" ]