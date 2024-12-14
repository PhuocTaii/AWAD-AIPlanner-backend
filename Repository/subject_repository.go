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
	var subject *models.Subject
	err := config.SubjectCollection.FindOne(ctx, bson.M{"_id": utils.ConvertStringToObjectID(id)}).Decode(&subject)
	if err != nil {
		return nil, err
	}
	return subject, nil
}

func FindAllUserSubject(ctx *gin.Context, userId string) ([]models.Subject, error) {
	cursor, err := config.SubjectCollection.Find(ctx, bson.M{"user._id": utils.ConvertStringToObjectID(userId), "is_deleted": false})
	if err != nil {
		return nil, err
	}
	var subjects []models.Subject
	if err = cursor.All(ctx, &subjects); err != nil {
		return nil, err
	}
	return subjects, nil
}

func InsertSubject(ctx *gin.Context, subject *models.Subject) (*models.Subject, error) {
	newSubject := &models.Subject{
		Name:      subject.Name,
		User:      subject.User,
		IsDeleted: false,
		CreatedAt: primitive.DateTime(utils.GetCurrentTime()),
		UpdatedAt: primitive.DateTime(utils.GetCurrentTime()),
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
