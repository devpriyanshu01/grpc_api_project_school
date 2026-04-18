package mongodb

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/devpriyanshu01/grpc_api_project_school/internals/models"
	"github.com/devpriyanshu01/grpc_api_project_school/pkg/utils"
	grpcapipb "github.com/devpriyanshu01/grpc_api_project_school/proto/gen"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddTeachersToDb(ctx context.Context, req *grpcapipb.Teachers) ([]*grpcapipb.Teacher, error) {
	//create mongo db client
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}
	defer client.Disconnect(ctx)

	log.Println("Inside Add Teachers service")

	newTeachers := make([]*models.Teacher, len(req.GetTeachers()))

	//all the sent teachers data is stored in newTeachers slice below
	j := 0
	for _, pbTeacher := range req.GetTeachers() {
		modelTeacher := models.Teacher{}
		pbVal := reflect.ValueOf(pbTeacher).Elem()        //reflect object of one sent teacher data
		modelVal := reflect.ValueOf(&modelTeacher).Elem() //reflect object of model

		for i := 0; i < pbVal.NumField(); i++ {
			pbField := pbVal.Field(i)               //gives the value of the field at i
			fieldName := pbVal.Type().Field(i).Name //name of field at i

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
	return addedTeachers, nil
}

func GetTeachersFromDb(ctx context.Context, req *grpcapipb.GetTeachersRequest) ([]*grpcapipb.Teacher, error) {
	log.Println("request inside GetTeachersDb")

	//filtering, getting the filters from the request, another function
	filter, err := BuildFilterForTeacher(req.Teacher, &models.Teacher{})
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error building filter object")
	}

	//sorting, getting the sort options from the request, another function
	sortOptions := buildSortOptions(req.GetSortBy())

	//connect to mongodb
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal server error")
	}
	defer client.Disconnect(ctx)

	coll := client.Database("school").Collection("teachers") //create collection instance
	var cursor *mongo.Cursor                                 //cursor is used to iterate over a stream of documents

	//if no sorting is required then we don't pass sortOptions to Find() fn.
	if len(sortOptions) < 1 {
		cursor, err = coll.Find(ctx, filter)
	} else {
		cursor, err = coll.Find(ctx, filter, options.Find().SetSort(sortOptions))
	}

	if err != nil {
		return nil, utils.ErrorHandler(err, "can't find teachers data")
	}
	defer cursor.Close(ctx)

	//iterate over each documents using Next function
	//we iterate over each documents that is streamed from the mongodb
	//The data that is received from the mongodb, can't be directly as protocol buffers.
	//Because mongodb returns bson & protocol buffers don't understand bson.
	//So, first we'll receive in go struct & then pass as protocol buffers.
	var teachers []*grpcapipb.Teacher
	for cursor.Next(ctx) {
		var teacher models.Teacher

		//cursor.Decode() can directly convert bson data to go struct because we've specified bson tags
		//while defining our models.Teacher struct. MongoDb will match each bson & assign the value to
		//teacher object wiz of type models.Teacher
		err = cursor.Decode(&teacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "decode from mongodb bson to go struct failed")
		}

		teachers = append(teachers, &grpcapipb.Teacher{
			Id:        teacher.Id,
			FirstName: teacher.FirstName,
			LastName:  teacher.LastName,
			Email:     teacher.Email,
			Class:     teacher.Class,
			Subject:   teacher.Subject,
		})
	}

	return teachers, nil
}

func buildFilterForTeacher(teacher *grpcapipb.Teacher) (bson.M, error) {
	filter := bson.M{}

	if teacher == nil {
		return filter, nil
	}

	var modelTeacher models.Teacher
	modelVal := reflect.ValueOf(&modelTeacher).Elem()
	modelType := modelVal.Type()

	//store the data from req to our internal struct
	reqVal := reflect.ValueOf(teacher).Elem()
	reqType := reqVal.Type()

	for i := 0; i < reqVal.NumField(); i++ {
		fieldVal := reqVal.Field(i)
		fieldName := reqType.Field(i).Name

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			modelField := modelVal.FieldByName(fieldName)
			if modelField.IsValid() && modelField.CanSet() {
				modelField.Set(fieldVal)
			}
		}
	}

	//Now we iterate over the modelTeacher to build using bson.M
	for i := 0; i < modelVal.NumField(); i++ {
		fieldVal := modelVal.Field(i)
		// fieldName := modelType.Field(i).Name

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			bsonTag := modelType.Field(i).Tag.Get("bson")
			bsonTag = strings.TrimSuffix(bsonTag, ",omitempty")
			
			//when request received is filter with _id the Mongodb won't be able to fetch the data.
			//Because _id is ObjectId is mongodb. So we need to convert the string-id received from 
			//the request to ObjectId.
			if bsonTag == "_id" {
				objectId, err := primitive.ObjectIDFromHex(teacher.Id)
				if err != nil {
					return nil, utils.ErrorHandler(err, "failed to convert given _id to objectId.")
				} 
				filter[bsonTag] = objectId
			}else {
				filter[bsonTag] = fieldVal.Interface().(string)
			}
		}
	}
	fmt.Println("Filter:", filter)

	return filter, nil
}

// for sorting
func buildSortOptions(SortFields []*grpcapipb.SortField) bson.D {
	var sortOptions bson.D

	for _, sortField := range SortFields {
		order := 1
		if sortField.GetOrder() == grpcapipb.Order_DESC {
			order = -1
		}
		sortOptions = append(sortOptions, bson.E{
			Key:   sortField.Field,
			Value: order,
		})
	}

	log.Println("SortOptions:", sortOptions)
	return sortOptions
}
