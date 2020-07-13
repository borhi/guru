package repositories

import (
	"context"
	"github.com/imdario/mergo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"guru/models"
	"time"
)

const TransactionCollection = "transaction"

type TransactionRepository struct {
	DB *mongo.Database
}

func (r *TransactionRepository) FindAllBet(statistic map[uint64]*models.StatisticModel) error {
	collection := r.DB.Collection(TransactionCollection)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	matchStage := bson.D{{"$match", bson.D{{"type", models.TypeBet}}}}
	groupStage := bson.D{{
		"$group",
		bson.D{
			{"_id", "$user_id"},
			{"bet_sum", bson.D{{"$sum", "$amount"}}},
			{"bet_count", bson.D{{"$sum", 1}}},
		},
	}}

	cur, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
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

func (r *TransactionRepository) FindAllWin(statistic map[uint64]*models.StatisticModel) error {
	collection := r.DB.Collection(TransactionCollection)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	matchStage := bson.D{{"$match", bson.D{{"type", models.TypeWin}}}}
	groupStage := bson.D{{
		"$group",
		bson.D{
			{"_id", "$user_id"},
			{"win_sum", bson.D{{"$sum", "$amount"}}},
			{"win_count", bson.D{{"$sum", 1}}},
		},
	}}

	cur, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
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

func (r *TransactionRepository) Insert(transactionModel models.TransactionModel) error {
	collection := r.DB.Collection(TransactionCollection)

	_, err := collection.InsertOne(context.TODO(), transactionModel)
	if err != nil {
		return err
	}

	return nil
}
