package args

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/octanegg/zsr/octane"
	"github.com/octanegg/zsr/octane/filter"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	debugInput  = flag.Bool("debug", false, "debug")
	teamsInput  = flag.String("team", "", "octane team names")
	groupsInput = flag.String("group", "", "octane group ids")
)

type Args struct {
	Debug  bool
	Teams  []string
	Groups []string
}

func BuildFilter(o octane.Client) (bson.M, error) {
	flag.Parse()

	var groups, teams []string

	if *groupsInput != "" {
		groups = strings.Split(*groupsInput, ",")
	}

	if *teamsInput != "" {
		teams = strings.Split(*teamsInput, ",")
	}

	for _, team := range teams {
		if _, err := o.Teams().FindOne(bson.M{"name": team}); err != nil {
			return nil, errors.New("team not found")
		}
	}

	f := filter.New(
		filter.Strings("team.team.name", teams),
		filter.Strings("game.match.event.groups", groups),
	)

	if *debugInput {
		dump, err := json.Marshal(f)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Using filter: %s\n", string(dump))
	}

	return f, nil
}
