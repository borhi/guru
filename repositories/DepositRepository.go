package repositories

import (
	"context"
	"github.com/imdario/mergo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"guru/models"
	"time"
)

const depositCollection = "deposit"

type DepositRepository struct {
	DB *mongo.Database
}

func (r *DepositRepository) FindAllDeposit(statistic map[uint64]*models.StatisticModel) error {
	collection := r.DB.Collection(depositCollection)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	groupStage := bson.D{{
		"$group",
		bson.D{
			{"_id", "$user_id"},
			{"deposit_sum", bson.D{{"$sum", "$amount"}}},
			{"deposit_count", bson.D{{"$sum", 1}}},
		},
	}}

	cur, err := collection.Aggregate(ctx, mongo.Pipeline{groupStage})
	if err != nil {
		return err
	}

	results := make(map[uint64]*models.StatisticModel)
	for cur.Next(ctx) {
		var result models.StatisticModel
		if err := cur.Decode(&result); err != nil {
			return err
		}

		results[result.Id] = &result
	}
	if err := cur.Err(); err != nil {
		return err
	}

	if err := mergo.Merge(&statistic, results); err != nil {
		return err
	}

	return nil
}

func (r *DepositRepository) Insert(depositModel models.DepositModel) error {
	collection := r.DB.Collection(depositCollection)

	_, err := collection.InsertOne(context.TODO(), depositModel)
	if err != nil {
		return err
	}

	return nil
}
