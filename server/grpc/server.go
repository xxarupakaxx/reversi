package main

import (
	"fmt"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	logrus2 "github.com/sirupsen/logrus"
	pb "github.com/xxarupakaxx/reversi/gen/pb/proto"
	"github.com/xxarupakaxx/reversi/server/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	port := 50051
	lis,err := net.Listen("tcp",fmt.Sprintf(":%d",port))
	if err != nil {
		log.Fatalf("failed to listen :%v",err)
	}
	logrusLogger := logrus2.New()
	logrusEnty := logrus2.NewEntry(logrusLogger)
	grpc_logrus.ReplaceGrpcLogger(logrusEnty)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_logrus.UnaryServerInterceptor(logrusEnty),
			),
		)

	pb.RegisterMatchingServiceServer(server,handler.NewMatchingHandler())
	pb.RegisterGameServiceServer(server,handler.NewGameHandler())

	reflection.Register(server)

	go func() {
		log.Println("start gRPC server port ", port)
		err = server.Serve(lis)
		if err != nil {
			log.Println("failed to start server :",err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit,os.Interrupt)
	<- quit
	log.Println("stopping gRPC server..")
	server.GracefulStop()
}
