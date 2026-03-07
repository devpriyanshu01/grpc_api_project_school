package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/devpriyanshu01/grpc_api_project_school/internals/api/handlers"
	"github.com/devpriyanshu01/grpc_api_project_school/internals/repositories/mongodb"
	pb "github.com/devpriyanshu01/grpc_api_project_school/proto/gen"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	_, err := mongodb.CreateMongoClient(context.Background())
	if err != nil {
		log.Println("Failed to connect to mongodb")
		return
	}
	

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return
	}

	s := grpc.NewServer()

	pb.RegisterStudentsServiceServer(s, &handlers.Server{})
	pb.RegisterTeachersServiceServer(s, &handlers.Server{})
	pb.RegisterExecsServiceServer(s, &handlers.Server{})

	reflection.Register(s)

	port := os.Getenv("SERVER_PORT")

	fmt.Println("gRPC Server is running on port:", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("Failed to create listener for gRPC server: ", err)
		return
	}

	err = s.Serve(lis)
	if err != nil {
		log.Fatal("Failed to serve: ", err)
	}
}
