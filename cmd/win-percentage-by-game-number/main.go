package main

import (
	"fmt"
	"os"

	"github.com/octanegg/stats/args"
	"github.com/octanegg/zsr/octane"
	"github.com/octanegg/zsr/octane/pipelines"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	o, err := octane.New(os.Getenv("DB"))
	if err != nil {
		panic(err)
	}

	filter, err := args.BuildFilter(o)
	if err != nil {
		fmt.Println(err)
		return
	}

	type game struct {
		Number int `bson:"_id"`
		Total  int `bson:"games"`
		Wins   int `bson:"wins"`
	}

	res, err := o.Statlines().Pipeline(
		pipelines.New(
			pipelines.Match(filter),
			pipelines.Group(bson.M{
				"_id": "$game.number",
				"games": bson.M{
					"$sum": bson.M{
						"$divide": bson.A{1, "$game.match.event.mode"},
					},
				},
				"wins": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{
								"$eq": bson.A{
									"$team.winner", true,
								},
							},
							bson.M{
								"$divide": bson.A{
									1,
									"$game.match.event.mode",
								},
							},
							0,
						},
					},
				},
			}),
			pipelines.Sort("_id", false),
		),
		func(cursor *mongo.Cursor) (interface{}, error) {
			var i game
			if err := cursor.Decode(&i); err != nil {
				return nil, err
			}
			return i, nil
		},
	)
	if err != nil {
		panic(err)
	}

	for _, r := range res {
		i := r.(game)
		fmt.Printf("Game %d: %.2f%% (%d games)\n", i.Number, float64(i.Wins)/float64(i.Total)*100, i.Total)
	}
}
