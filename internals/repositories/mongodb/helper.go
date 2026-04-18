package mongodb

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/devpriyanshu01/grpc_api_project_school/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// let's make this function generalized using interfaces so that it can work for teacher, students, execs or anything.
func BuildFilterForTeacher(object interface{}, model interface{}) (bson.M, error) {
	filter := bson.M{}

	if object == nil || reflect.ValueOf(object).IsNil() {
		return filter, nil
	}

	// var modelTeacher models.Teacher
	modelVal := reflect.ValueOf(model).Elem()
	modelType := modelVal.Type()

	//store the data from req to our internal struct
	reqVal := reflect.ValueOf(object).Elem()
	reqType := reqVal.Type()

	//copying request object to model object
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
		fieldName := modelType.Field(i).Name

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			bsonTag := modelType.Field(i).Tag.Get("bson")
			bsonTag = strings.TrimSuffix(bsonTag, ",omitempty")

			//when request received is filter with _id the Mongodb won't be able to fetch the data.
			//Because _id is ObjectId is mongodb. So we need to convert the string-id received from
			//the request to ObjectId.
			if bsonTag == "_id" {
				objectId, err := primitive.ObjectIDFromHex(fieldVal.FieldByName(fieldName).Interface().(string))
				if err != nil {
					return nil, utils.ErrorHandler(err, "failed to convert given _id to objectId.")
				}
				filter[bsonTag] = objectId
			} else {
				filter[bsonTag] = fieldVal.Interface().(string)
			}
		}
	}
	fmt.Println("Filter:", filter)

	return filter, nil
}
