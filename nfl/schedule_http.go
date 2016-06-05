package nfl

import (
    "net/http"
    "fmt"
    "io/ioutil"
    "encoding/xml"
)

func GetFullSeasonScheduleAsync(season int, c chan <- SeasonSchedule) {
    c <- GetFullSeasonSchedule(season)
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
        schedule.Weeks[i].WeekInit()
    }

    if len(schedule.Weeks) > 0 {
        return schedule.Weeks[0]
    } else {
        return Week{}
    }
}