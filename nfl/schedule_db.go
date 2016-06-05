package nfl

import (
    "log"
    "github.com/gocraft/dbr"
    "time"
    "fmt"
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
        var eid string

        weekAndMap := dbr.AndMap{
            "season" : week.Season,
            "season_type": seasonType,
            "week_num":    week.Week,
            "home_team":   g.HomeTeam,
            "away_team":   g.AwayTeam,
        }


        //log.Printf("EID: %s, Game Time: %s, Game Date Time: %s",g.Eid, g.Time, g.DateTime.Format(time.RFC3339))
        nflDb.Select("eid").From("nfl_schedule").Where(weekAndMap).Load(&eid)

        if eid != "" {
            up := nflDb.Update("nfl_schedule").
            Set("game_time", g.Time ).
            Set("game_datetime", g.DateTime.Format(time.RFC3339)).
            Set("home_score", g.HomeScore ).
            Set("away_score", g.AwayScore ).
            Set("last_updated_datetime",dbr.Now).
            Where(weekAndMap)
             _, err := up.Exec()
            if err != nil {
                log.Printf("error inserting into nfl_schedule, season: %v week: %v home: %s away: %s", week.Season, week.Week, g.HomeTeam, g.AwayTeam)
                log.Print(err)
                tran.Rollback()
                break
            }

        } else {
            g.SeasonType = seasonType
            g.Season = week.Season
            g.WeekNumber = week.Week

            //preferrably we could do this, but timezone info is not being passed in.
            //ins := nflDb.InsertInto("nfl_schedule").Columns("eid", "gsis","season","season_type","week_num","game_time","game_datetime","home_team","home_score", "away_team", "away_score").Record(&g)


            //Timestamp is not properly propograting timezone information into postgres, therefore we need to manually build our insert statement.
            ins := nflDb.InsertBySql("insert into nfl_schedule (eid, gsis, season, season_type, week_num, game_time, game_datetime, home_team, home_score, away_team, away_score) values (" +
                fmt.Sprintf("'%s', %v, %v, '%s', %v, '%s', '%s', '%s', %v, '%s', %v",g.Eid, g.Gsis, g.Season, g.SeasonType, g.WeekNumber, g.Time, g.DateTime.Format(time.RFC3339), g.HomeTeam, g.HomeScore, g.AwayTeam, g.AwayScore) + ")")
            //log.Print(ins2.ToSql())
            _, err := ins.Exec()
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
    nfldb.Select("*").From("nfl_schedule").Where(weekExistsAndMap).Load(&games)
    week.Games = games
    return week;

}