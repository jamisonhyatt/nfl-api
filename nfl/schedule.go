package nfl

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"golang.org/x/net/context"
	"github.com/guregu/db"
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
	GamesDictionary map[string]Game `json:"-"`
}

func GetFullSeasonSchedule(season int) SeasonSchedule {
	fullSchedule := SeasonSchedule{Season: season}
	fullSchedule.PreSeason = GetSeasonTypeSchedule(season, "PRE")
	fullSchedule.RegularSeason = GetSeasonTypeSchedule(season, "REG")
	fullSchedule.PostSeason = GetSeasonTypeSchedule(season, "POST")

	return fullSchedule
}

func GetSeasonTypeSchedule(season int, seasonType string) Schedule {
	var weekRangeStart int
	var weekRangeEnd int
	if seasonType == "POST" {
		weekRangeStart = 18
		weekRangeEnd = 22
	} else if seasonType == "PRE" {
		weekRangeStart = 0
		weekRangeEnd = 4
	} else {
		weekRangeStart = 1
		weekRangeEnd = 17
	}

	schedule := Schedule{SeasonType: seasonType}

	for i := weekRangeStart; i <= weekRangeEnd; i++ {
		if i == 21 {
			continue
		} //WTF NFL...
		week := GetWeeklySchedule(season, seasonType, i)
		if len(week.Games) > 0 {
			schedule.Weeks = append(schedule.Weeks, week)
		}
	}
	if len(schedule.Weeks) == 0 {
		schedule.Weeks = make(Weeks, 0, 0)
	}
	return schedule

}

func GetWeeklySchedule(season int, seasonType string, week int) Week {

	req, _ := http.NewRequest("GET", "http://www.nfl.com/ajax/scorestrip", nil)
	query := req.URL.Query()
	query.Add("season", fmt.Sprintf("%d", season))
	query.Add("seasonType", seasonType)
	query.Add("week", fmt.Sprintf("%d", week))

	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	xmlSchedule, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var schedule Schedule

	xml.Unmarshal(xmlSchedule, &schedule)

	for i, _ := range schedule.Weeks {
		schedule.Weeks[i].Init()
	}

	if len(schedule.Weeks) > 0 {
		return schedule.Weeks[0]
	} else {
		return Week{}
	}
}

func (w *Week) Init() {
	w.GamesDictionary = make(map[string]Game, len(w.Games))

	saturdayGames := make([]Game, 0, len(w.Games))
	sundayGames := make([]Game, 0, len(w.Games))

	for i, game := range w.Games {
		switch game.Day {
		case "Sun":
			sundayGames = append(sundayGames, game)
		case "Sat":
			saturdayGames = append(saturdayGames, game)
		default:
			{
				w.Games[i].Period = "PM"
				w.Games[i].DateTime = w.Games[i].JSDateTime()
				w.GamesDictionary[game.Eid] = w.Games[i]
			}
		}
	}

	if len(saturdayGames) > 0 {
		SetPeriod(saturdayGames)
		for i, _ := range saturdayGames {
			saturdayGames[i].DateTime = saturdayGames[i].JSDateTime()
			w.GamesDictionary[saturdayGames[i].Eid] = saturdayGames[i]
		}
	}

	SetPeriod(sundayGames)
	for i, _ := range sundayGames {
		sundayGames[i].DateTime = sundayGames[i].JSDateTime()
		w.GamesDictionary[sundayGames[i].Eid] = sundayGames[i]
	}

	//mirror our list and
	var gi = 0
	for _, game := range w.GamesDictionary {
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

 func AddGameToDb (w Week, ctx context.Context, wg sync.WaitGroup)  {
     defer wg.Done()

     nflDb := db.SQL(ctx, "nfl")
	 tx, _ := nflDb.Begin()

     stmt, err := tx.Prepare("call InsertUpdateSchedule_p (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
     defer stmt.Close()
     if err != nil {
         log.Fatal(err)
     }

	 for _, game := range(w.Games) {
		 stmt.Exec(game.Eid, game.Gsis, w.)
	 }

	 //CREATE PROCEDURE InsertUpdateSchedule_p
	 //( IN in_eid varchar(10)
	 //,IN in_gsis int
	 //,IN in_season_type char(4)
	 //,IN in_week_num   int
	 //,IN in_game_time  varchar(5)
	 //,IN in_game_datetime datetime
	 //,IN in_home_team char(30)
	 //,IN in_home_score int
	 //,IN in_away_team char(3)
	 //,IN in_away_score int
	 //)

     stmt.Exec(g.Eid, g.Gsis, )

 }
