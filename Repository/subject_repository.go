package repository

import (
	config "project/Config"
	models "project/Models"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindSubjectById(ctx *gin.Context, id string) (*models.Subject, error) {
	filter := bson.M{
		"_id":        utils.ConvertStringToObjectID(id),
		"is_deleted": false,
	}
	var subject *models.Subject
	err := config.SubjectCollection.FindOne(ctx, filter).Decode(&subject)
	if err != nil {
		return nil, err
	}
	return subject, nil
}

func FindSubjectByIdAndUserId(ctx *gin.Context, id, userId string) (*models.Subject, error) {
	filter := bson.M{
		"_id":        utils.ConvertStringToObjectID(id),
		"user":       utils.ConvertStringToObjectID(userId),
		"is_deleted": false,
	}
	var subject *models.Subject
	err := config.SubjectCollection.FindOne(ctx, filter).Decode(&subject)
	if err != nil {
		return nil, err
	}
	return subject, nil
}

func FindAllUserSubject(ctx *gin.Context, userId string) ([]models.Subject, error) {
	filter := bson.M{
		"user":       utils.ConvertStringToObjectID(userId),
		"is_deleted": false,
	}
	cursor, err := config.SubjectCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var subjects []models.Subject
	if err = cursor.All(ctx, &subjects); err != nil {
		return nil, err
	}
	return subjects, nil
}

func IsSubjectExisted(ctx *gin.Context, name, userId string) bool {
	filter := bson.M{
		"name":       name,
		"user":       utils.ConvertStringToObjectID(userId),
		"is_deleted": false,
	}
	var subject *models.Subject
	isExist := config.SubjectCollection.FindOne(ctx, filter).Decode(&subject)
	return isExist != nil
}

func InsertSubject(ctx *gin.Context, subject *models.Subject) (*models.Subject, error) {
	newSubject := &models.Subject{
		Name:      subject.Name,
		User:      subject.User,
		IsDeleted: false,
		CreatedAt: utils.GetCurrent(),
		UpdatedAt: utils.GetCurrent(),
	}

	res, err := config.SubjectCollection.InsertOne(ctx, newSubject)
	if err != nil {
		return nil, err
	}

	response := &models.Subject{
		ID:        res.InsertedID.(primitive.ObjectID),
		Name:      newSubject.Name,
		User:      newSubject.User,
		IsDeleted: false,
		CreatedAt: newSubject.CreatedAt,
		UpdatedAt: newSubject.UpdatedAt,
	}

	return response, nil
}

func UpdateSubject(ctx *gin.Context, subject *models.Subject) (*models.Subject, error) {
	filter := bson.M{
		"_id":        subject.ID,
		"is_deleted": false,
	}
	update := bson.M{"$set": bson.M{
		"name":       subject.Name,
		"updated_at": primitive.DateTime(utils.GetCurrentTime()),
	}}

	_, err := config.SubjectCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return subject, nil
}

func DeleteSubject(ctx *gin.Context, subject *models.Subject) (*models.Subject, error) {
	filter := bson.M{
		"_id":        subject.ID,
		"is_deleted": false,
	}
	update := bson.M{"$set": bson.M{
		"is_deleted": true,
		"updated_at": primitive.DateTime(utils.GetCurrentTime()),
	}}
	_, err := config.SubjectCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return subject, nil
}
