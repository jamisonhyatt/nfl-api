package nfl

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/guregu/db"
	"golang.org/x/net/context"
)

type Conference struct {
	Name      string     `json:"name"`
	Divisions []Division `json:"divisions"`
}

type Division struct {
	Region         string `json:"region"`
	ConferenceName string `json:"conferenceName"`
	Teams          []Team `json:"teams"`
}

type Teams []Team
type Team struct {
	TeamId         string `json:"teamId"`
	NickName       string `json:"nickName"`
	City           string `json:"city"`
	DivisionName   string `json:"divisionName"`
	ConferenceName string `json:"conferenceName"`
}

func GetTeam(teamId string, ctx context.Context, c chan Teams) {
	getTeams(teamId, ctx, c)
}

func GetAllTeams(ctx context.Context, c chan Teams) {
	getTeams("", ctx, c)
}

//returns a channel of Teams; All Teams are r
func getTeams(teamID string, ctx context.Context, c chan Teams) {
	var teams Teams

	nflDb := db.SQL(ctx, "nfl")

	stmt, err := nflDb.Prepare("call GetTeam_p (?)")
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	var rows *sql.Rows

	if len(strings.TrimSpace(teamID)) == 0 {
		rows, err = stmt.Query(nil)
	} else {
		rows, err = stmt.Query(teamID)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		teams = append(teams, ScanTeam(rows))
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	c <- teams
}

func ScanTeam(r *sql.Rows) Team {
	var team Team
	r.Scan(&team.TeamId, &team.City, &team.NickName, &team.DivisionName, &team.ConferenceName)
	return team
}
