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

	a, err := args.Get(o)
	if err != nil {
		panic(err)
	}

	type appearances struct {
		Player      *octane.Player `bson:"player"`
		Team        *octane.Team   `bson:"team"`
		Appearances int            `bson:"appearances"`
	}

	res, err := o.Statlines().Pipeline(
		pipelines.New(
			pipelines.Match(bson.M{"game.match.event.mode": 3}),
			pipelines.Group(bson.M{
				"_id": bson.M{
					"player": "$player.player._id",
					"team":   "$team.team._id",
				},
				"player": bson.M{
					"$first": "$player.player",
				},
				"team": bson.M{
					"$first": "$team.team",
				},
				"appearances": bson.M{
					"$sum": 1,
				},
			}),
			pipelines.Sort("appearances", true),
			pipelines.Limit(a.Limit),
		),
		func(cursor *mongo.Cursor) (interface{}, error) {
			var i appearances
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
		i := r.(appearances)
		fmt.Printf("%4d %-20s %-10s\n", i.Appearances, i.Team.Name, i.Player.Tag)
	}
}
