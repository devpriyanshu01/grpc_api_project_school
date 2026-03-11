package handlers

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/devpriyanshu01/grpc_api_project_school/internals/models"
	"github.com/devpriyanshu01/grpc_api_project_school/internals/repositories/mongodb"
	"github.com/devpriyanshu01/grpc_api_project_school/pkg/utils"
	grpcapipb "github.com/devpriyanshu01/grpc_api_project_school/proto/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Server) AddTeachers(ctx context.Context, req *grpcapipb.Teachers) (*grpcapipb.Teachers, error) {
	//create mongo db client
	client, err := mongodb.CreateMongoClient(ctx)
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}
	defer client.Disconnect(ctx)

	log.Println("INside Add Teachers service")

	fmt.Println(req.GetTeachers())

	newTeachers := make([]*models.Teacher, len(req.GetTeachers()))

	for _, pbTeacher := range req.GetTeachers() {
		modelTeacher := models.Teacher{FirstName: "Raman"}
		pbVal := reflect.ValueOf(pbTeacher).Elem() //reflect object of one sent teacher data
		modelVal := reflect.ValueOf(&modelTeacher).Elem() //reflect object of mode

		j := 0
		for i := 0; i < pbVal.NumField(); i++ {
			pbField := pbVal.Field(i)	//gives the value of the field at i
			fieldName := pbVal.Type().Field(i).Name  //name of field at i

			modelField := modelVal.FieldByName(fieldName)
			if modelField.IsValid() && modelField.CanSet() {
				modelField.Set(pbField)
			}
		}

		newTeachers[j] = &modelTeacher
		j++
	}

	var addedTeachers []*grpcapipb.Teacher //protobuf struct

	for _, teacher := range newTeachers {
		result, err := client.Database("school").Collection("teachers").InsertOne(ctx, *teacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error adding value to the database")
		}

		objectId, ok := result.InsertedID.(primitive.ObjectID)
		if ok {
			teacher.Id = objectId.Hex()
		}

		//send the newTeachers value to protobuf Teachers struct
		pbTeacher := &grpcapipb.Teacher{}
		modelVal := reflect.ValueOf(*teacher)
		pbVal := reflect.ValueOf(pbTeacher).Elem()

		for i := 0; i < modelVal.NumField(); i++ {
			fieldVal := modelVal.Field(i)
			fieldName := modelVal.Type().Field(i).Name

			pbField := pbVal.FieldByName(fieldName)

			if pbField.IsValid() && pbField.CanSet() {
				pbField.Set(fieldVal)
			}
		}
		addedTeachers = append(addedTeachers, pbTeacher)

	}
	return &grpcapipb.Teachers{Teachers: addedTeachers}, nil
}