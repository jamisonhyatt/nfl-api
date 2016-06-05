package nfl

import (
    "fmt"
    "log"
    "strconv"
    "strings"
    "time"
)

//Games is a slice of Game structs
type Games []Game

//Game that represents a single game within a season schedule
type Game struct {
    Eid       string `xml:"eid,attr" json:"eid"`
    Gsis      int    `xml:"gsis,attr" json:"gsis"`
    Season      int     `db:"season"`
    SeasonType string `db:"season_type"`
    WeekNumber int `db:"week_num"`
    Day       string `xml:"d,attr" json:"Day"`
    Time      string `xml:"t,attr" json:"gameTime" db:"game_time"`
    Quarter   string `xml:"q,attr" json:"quarter"`
    HomeTeam  string `xml:"h,attr" json:"homeTeam" db:"home_team"`
    HomeScore int    `xml:"hs,attr" json:"homeScore" db:"home_score"`
    AwayTeam  string `xml:"v,attr" json:"awayTeam" db:"away_team"`
    AwayScore int    `xml:"vs,attr" json:"awayScore" db:"away_score"`

    Period   string    `json:"period"` //AM PM
    DateTime time.Time `json:"gameDateTime" db:"game_datetime"`
}

type ByGameSequence Games

func (a ByGameSequence) Len() int           { return len(a) }
func (a ByGameSequence) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByGameSequence) Less(i, j int) bool { return a[i].GameDaySequence() < a[j].GameDaySequence() }

type ByEid Games

func (a ByEid) Len() int           { return len(a) }
func (a ByEid) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByEid) Less(i, j int) bool { return a[i].Eid < a[j].Eid }

func (g Game) Year() string {
    return g.Eid[:4]
}

func (g Game) Month() string {
    return g.Eid[4:6]
}

func (g Game) DayMonth() string {
    return g.Eid[6:8]
}

func (g Game) GameDaySequence() int {
    i, _ := strconv.Atoi(g.Eid[8:10])
    return i
}

func (g Game) Hour() int {
    t := strings.Split(g.Time, ":")
    i, _ := strconv.Atoi(t[0])
    return i
}

func (g Game) Minute() int {
    t := strings.Split(g.Time, ":")
    i, _ := strconv.Atoi(t[2])
    return i
}

func JSTime(g Game) string {

    s := strings.Split(g.Time, ":")
    hour, _ := strconv.Atoi(s[0])

    if g.Period == "PM" && hour >= 1 && hour < 12 {
        hour += 12
    } else if g.Period == "AM" && hour < 1 {
        hour = 0
    }
    return fmt.Sprintf("T%s:%s:00", Right("0"+strconv.Itoa(hour), 2), s[1])
}

//

func (g Game) JSDateTime() time.Time {
    var loc, _ = time.LoadLocation("America/New_York") //All Game times EST
    const RFC3339NoZone = "2006-01-02T15:04:05"

    jstime := JSTime(g)

    dttmString := fmt.Sprintf("%s-%s-%s%s", g.Year(), g.Month(), g.DayMonth(), jstime)

    t, err := time.ParseInLocation(RFC3339NoZone, dttmString, loc)
    if err != nil {
        log.Fatal(err)
    }
    return t

}
