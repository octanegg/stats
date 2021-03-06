package args

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/octanegg/zsr/octane"
	"github.com/octanegg/zsr/octane/filter"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	debugInput  = flag.Bool("debug", false, "debug")
	teamsInput  = flag.String("team", "", "octane team names")
	groupsInput = flag.String("group", "", "octane group ids")
	limitInput  = flag.String("limit", "", "limit records")
)

type Args struct {
	Debug  bool
	Teams  []string
	Groups []string
	Limit  int
}

func Get(o octane.Client) (*Args, error) {
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

	limit, _ := strconv.Atoi(*limitInput)

	return &Args{
		Teams:  teams,
		Groups: groups,
		Limit:  limit,
	}, nil
}

func BuildFilter(o octane.Client) (bson.M, error) {
	args, err := Get(o)
	if err != nil {
		return nil, err
	}

	f := filter.New(
		filter.Strings("team.team.name", args.Teams),
		filter.Strings("game.match.event.groups", args.Groups),
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
