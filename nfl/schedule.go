package nfl

import (
    "sort"
    "github.com/gocraft/dbr"
    "sync"
    "time"
    "encoding/xml"
    "log"
)

type SeasonSchedule struct {
    Season        int      `json:"season"`
    PreSeason     Schedule `json:"preSeason"`
    RegularSeason Schedule `json:"regularSeason"`
    PostSeason    Schedule `json:"postSeason`
}

type Schedule struct {
    XMLName    xml.Name `xml:"ss" json:"-"`
    SeasonType string   `json:"seasonType"`
    Weeks      Weeks    `xml:"gms"`
}

type Weeks []Week
type Week struct {
    XMLName         xml.Name        `xml:"gms" json:"-"`
    Week            int             `xml:"w,attr" json:"week"`
    Season          int             `xml:"y,attr" json:"season"`
    Games           Games           `xml:"g" json:"games"`
}

//last nfl reorg, just keeping it simple.
var firstSeason = 2002

func InitDbSchedule(conn *dbr.Connection) {
    firstWeekInDb := DbGetWeeklySchedule(conn, "REG",firstSeason,  1)
    if (len(firstWeekInDb.Games) == 0) {
        log.Printf("no games in first season (%v), backfilling", firstSeason)
        start := time.Now()
        BackFillSeasonsAsync(conn)
        elapsed := time.Since(start)

        log.Printf("Backfill tooktook %s\n", elapsed)

    }
}


func (w *Week) WeekInit() {
    gamesDictionary := make(map[string]Game, len(w.Games))

    saturdayGames := make([]Game, 0, len(w.Games))
    sundayGames := make([]Game, 0, len(w.Games))
    season := w.Season
    for i, game := range w.Games {
        w.Games[i].Season = w.Season
        switch game.Day {
        case "Sun":
            sundayGames = append(sundayGames, game)
        case "Sat":
            saturdayGames = append(saturdayGames, game)
        default:
            {
                w.Games[i].Period = "PM"
                w.Games[i].DateTime = w.Games[i].JSDateTime()
                gamesDictionary[game.Eid] = w.Games[i]
            }
        }

    }

    if len(saturdayGames) > 0 {
        SetPeriod(saturdayGames)
        for i, _ := range saturdayGames {
            saturdayGames[i].DateTime = saturdayGames[i].JSDateTime()
            gamesDictionary[saturdayGames[i].Eid] = saturdayGames[i]
        }
    }

    SetPeriod(sundayGames)
    for i, _ := range sundayGames {
        sundayGames[i].DateTime = sundayGames[i].JSDateTime()
        gamesDictionary[sundayGames[i].Eid] = sundayGames[i]
    }

    //mirror our list and cleanup data
    var gi = 0
    for _, game := range gamesDictionary {
        if season < 2016 {
            if game.HomeTeam == "JAC" {
                game.HomeTeam = "JAX"
            } else if game.HomeTeam == "STL" {
                game.HomeTeam = "LA"
            }

            if game.AwayTeam == "JAC" {
                game.AwayTeam = "JAX"
            } else if game.AwayTeam == "STL" {
                game.AwayTeam = "LA"
            }

        }
        w.Games[gi] = game
        gi++
    }
    sort.Sort(ByEid(w.Games))
}

func SetPeriod(games Games) {
    sort.Sort(ByGameSequence(games))

    for i, game := range games {
        //If it's before 6 EST, assume PM.
        if game.Hour() < 6 {
            games[i].Period = "PM"
            continue
        }
        //Optimistically, Last game on a Sat or Sunday has gotta be in the afternoon.
        if i == len(games)-1 {
            games[i].Period = "PM"
            continue
        }
        var next Game
        next = games[i+1]
        if game.Hour() > next.Hour() && game.GameDaySequence() < next.GameDaySequence() {
            games[i].Period = "AM"
            continue
        }
        games[i].Period = "PM"
    }
}


func BackFillSeasonsAsync(conn *dbr.Connection) {

    var wg sync.WaitGroup
    seasonEnd := time.Now().Year()+1
    ch := make(chan SeasonSchedule, seasonEnd - firstSeason)
    for i := firstSeason; i < seasonEnd; i++ {
        wg.Add(1)
        go GetFullSeasonScheduleAsync(i, ch)
    }


    go func() {
        for sched := range ch {
            WriteFullSeasonSchedule(conn, sched)
            wg.Done()
        }
    }()
    wg.Wait()

}
