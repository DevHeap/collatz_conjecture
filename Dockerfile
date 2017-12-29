FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get github.com/gorilla/websocket
RUN go build src/backend/main.go
EXPOSE 80
ENTRYPOINT ["/app/main"]