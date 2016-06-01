package main

import (
	"encoding/json"
	"fmt"
	"nfl-api/nfl"
	"os"
	"time"

	"github.com/guregu/db"
	"golang.org/x/net/context"
	_ "github.com/lib/pq"
)

var jsonIndent = "  "

func main() {
	TestDBConnection()
	TestAPIScheduleRequest()

}

func TestAPIScheduleRequest() {
	sched := nfl.GetFullSeasonSchedule(2015)
	jsonStr, _ := json.MarshalIndent(sched, "", jsonIndent)
	os.Stdout.Write(jsonStr)
}

func TestDBConnection() {
	ctx := context.Background()
	//db, err := sql.Open("postgres", "host=192.168.2.101 port=5432 dbname=nfl user=nfl_api  password=nfl_api sslmode=require")

	//ctx = db.OpenSQL(ctx, "nfl", "mysql", "nfl_api:nfl_pass@tcp(192.168.2.101:3306)/nfl?autocommit=true")
	//db, err := sql.Open("postgres", "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")

	//ctx = db.OpenSQL(ctx, "nfl", "mysql", "root:root@tcp(172.16.102.129:3306)/nfl?autocommit=true")
	ctx = db.OpenSQL(ctx, "nfl", "postgres", "host=192.168.2.101 port=5432 dbname=nfl user=nfl_api  password=nfl_api sslmode=disable")
	defer db.Close(ctx) // closes all DB connections
	teamChannel := make(chan nfl.Teams)
	go nfl.GetTeam("SEA", ctx, teamChannel)
	seattle := <-teamChannel

	jsonStr, _ := json.MarshalIndent(seattle, "", jsonIndent)
	os.Stdout.Write(jsonStr)

	start := time.Now()
	iter := 2
	teamChannel = make(chan nfl.Teams)
	teamBlock := make([]nfl.Teams, iter, iter)

	for i := 0; i < iter; i++ {
		go nfl.GetAllTeams(ctx, teamChannel)
		teams := <-teamChannel
		teamBlock[i] = teams
	}

	elapsed := time.Since(start)

	jsonStr, _ = json.MarshalIndent(teamBlock[0], "", jsonIndent)
	os.Stdout.Write(jsonStr)
	fmt.Printf("\n%v database proc executions took %s\n", iter, elapsed)

	teamTotal := 0
	for i := 0; i < len(teamBlock); i++ {
		teamTotal += len(teamBlock[i])
	}
	fmt.Printf("rows retreived: %v", teamTotal)

}

//func printSchedule(season int) {
//    x := nfl.GetFullSeasonSchedule(season)
//
//    b, err := json.MarshalIndent(x, "", "\t")
//    if err != nil {
//        fmt.Println("error:", err)
//    }
//    os.Stdout.Write(b)
//}

//func readFile(s string) []byte {
//    xmlFile, err := os.Open(s)
//    var b []byte
//    if err != nil {
//        fmt.Println("Error opening file:", err)
//        return b
//    }
//    defer xmlFile.Close()
//
//    b, _ = ioutil.ReadAll(xmlFile)
//
//    return b;
//}
//
//func processXml(b []byte) {
//
//
//    var schedule nfl.Schedule
//
//    xml.Unmarshal(b, &schedule)
//
//    for i, _ := range(schedule.Weeks) {
//        schedule.Weeks[i].Init();
//
//        fmt.Println("gameDictionaryCount = ", len(schedule.Weeks[i].GamesDictionary))
//
//        //for _, game := range (schedule.Weeks[i].GamesDictionary) {
//        //    fmt.Printf("%s : %s-%s-%s, %d, %s\n", game.Eid, game.Year(), game.Month(), game.DayMonth(), game.GameDaySequence(), game.Period)
//        //}
//
//        for _, game := range (schedule.Weeks[i].Games) {
//            fmt.Printf("%s : %v, %d, %s\n", game.Eid, game.JSDateTime(), game.GameDaySequence(), game.Period)
//        }
//    }
//}
