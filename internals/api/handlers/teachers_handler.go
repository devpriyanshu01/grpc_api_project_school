package handlers

import (
	"context"
	"fmt"

	"github.com/devpriyanshu01/grpc_api_project_school/internals/repositories/mongodb"
	grpcapipb "github.com/devpriyanshu01/grpc_api_project_school/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AddTeachers(ctx context.Context, req *grpcapipb.Teachers) (*grpcapipb.Teachers, error) {
	addedTeachers, err := mongodb.AddTeachersToDb(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &grpcapipb.Teachers{Teachers: addedTeachers}, nil
}

func (s *Server) GetTeachers(ctx context.Context, req *grpcapipb.GetTeachersRequest) (*grpcapipb.Teachers, error) {
	teachers, err := mongodb.GetTeachersFromDb(ctx, req)
	if err != nil {
		fmt.Println("teachers;", teachers)
	}
	//dummy return
	return nil, nil
}
