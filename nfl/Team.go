package nfl

import (
    "database/sql"

    "github.com/gocraft/dbr"
    _ "github.com/lib/pq"
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
    TeamId         string `json:"teamId" db:"team_id"`
    NickName       string `json:"nickName" db:"nickname"`
    City           string `json:"city"`
    DivisionName   string `json:"divisionName" db:"division"`
    ConferenceName string `json:"conferenceName" db:"conference"`
}

func GetTeam(teamId string, conn *dbr.Connection, c chan Teams) {
    getTeams(teamId, conn, c)
}

func GetAllTeams(conn *dbr.Connection, c chan Teams) {
    getTeams("", conn, c)
}

//returns a channel of Teams; All Teams are r
func getTeams(teamID string, conn *dbr.Connection, c chan Teams) {
    var teams Teams

    nflDb := conn.NewSession(nil)
    if teamID == "" {
        nflDb.Select("*").From("team").Load(&teams)
    } else {
        nflDb.Select("*").From("team").Where("team_id = ?", teamID).Load(&teams)
    }
    c <- teams
}

func ScanTeam(r *sql.Rows) Team {
    var team Team
    r.Scan(&team.TeamId, &team.City, &team.NickName, &team.DivisionName, &team.ConferenceName)
    return team
}
