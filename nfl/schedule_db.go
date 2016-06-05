package nfl

import (
    "log"
    "github.com/gocraft/dbr"
)

func WriteFullSeasonSchedule(conn *dbr.Connection, fullSeason SeasonSchedule) {
    WriteSeasonSchedule(conn, fullSeason.PreSeason)
    WriteSeasonSchedule(conn, fullSeason.RegularSeason)
    WriteSeasonSchedule(conn, fullSeason.PostSeason)
}

func WriteSeasonSchedule(conn *dbr.Connection, schedule Schedule) {

    for _, w := range(schedule.Weeks){
        if len(w.Games) > 0 {
            WriteWeeklySchedule(conn, w, schedule.SeasonType)
        }
    }
}

func WriteWeeklySchedule(conn *dbr.Connection, week Week, seasonType string) {

    nflDb := conn.NewSession(nil)
    tran, err := nflDb.Begin()

    if err != nil {
        log.Print(err)
    }
    defer tran.RollbackUnlessCommitted()

    for _, g := range(week.Games) {
        var id int

        weekAndMap := dbr.AndMap{
            "season_type": seasonType,
            "week_num":    week.Week,
            "home_team":   g.HomeTeam,
            "away_team":   g.AwayTeam,
        }

        nflDb.Select("id").From("nfl_schedule").Where(weekAndMap).Load(&id)

        if id != 0 {
            nflDb.Update("nfl_schedule").
            Set("game_time", g.Time ).
            Set("game_datetime", g.DateTime ).
            Set("home_score", g.HomeScore ).
            Set("away_score", g.AwayScore ).
            Set("last_updated_datetime",dbr.Now).
            Where("id", id)

        } else {
            g.SeasonType = seasonType
            g.Season = week.Season
            g.WeekNumber = week.Week
            x := nflDb.InsertInto("nfl_schedule").Columns("eid", "gsis","season","season_type","week_num","game_time","game_datetime","home_team","home_score", "away_team", "away_score").Record(g)
            _, err := x.Exec()
            if err != nil {
                log.Printf("error inserting into nfl_schedule, season: %v week: %v home: %s away: %s", week.Season, week.Week, g.HomeTeam, g.AwayTeam)
                log.Print(err)
                tran.Rollback()
                break
            }
        }
    }
    if err != nil {
        err = tran.Commit()
    }

    if err != nil {
        log.Print(err)
    }
}

func DbGetWeeklySchedule(conn *dbr.Connection, seasonType string, season int, weekNum int) (Week){
    nfldb := conn.NewSession(nil)
    week := Week {Season : season, Week: weekNum, }

    weekExistsAndMap := dbr.AndMap{
        "season" : season,
        "season_type": seasonType,
        "week_num":    weekNum,
    }

    var games Games
    nfldb.Select("id").From("nfl_schedule").Where(weekExistsAndMap).Load(&games)
    week.Games = games
    return week;

}