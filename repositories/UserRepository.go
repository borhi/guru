package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"guru/models"
	"time"
)

const userCollection = "user"

type UserRepository struct {
	DB *mongo.Database
}

func (r *UserRepository) FindAll(users map[uint64]*models.UserModel) error {
	collection := r.DB.Collection(userCollection)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result models.UserModel
		if err := cur.Decode(&result); err != nil {
			return err
		}
		users[result.Id] = &result
	}
	if err := cur.Err(); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Insert(users []interface{}) error {
	collection := r.DB.Collection(userCollection)
	_, err := collection.InsertMany(context.TODO(), users)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(user *models.UserModel) error {
	collection := r.DB.Collection(userCollection)
	filter := bson.D{{"id", user.Id}}

	_, err := collection.UpdateOne(context.TODO(), filter, user)
	if err != nil {
		return err
	}

	return nil
}
