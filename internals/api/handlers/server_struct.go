package handlers

import pb "github.com/devpriyanshu01/grpc_api_project_school/proto/gen"

type Server struct {
	pb.UnimplementedStudentsServiceServer
	pb.UnimplementedTeachersServiceServer
	pb.UnimplementedExecsServiceServer
}