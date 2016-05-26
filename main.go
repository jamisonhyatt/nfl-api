package main

import (
    "nfl-api/nfl"
    "sync"
    "golang.org/x/net/context"
    "github.com/guregu/db"
    "fmt"
    "time"
    "encoding/json"
    "os"
    "log"
)


func main() {

    ctx := context.Background()
    ctx = db.OpenSQL(ctx, "nfl", "mysql", "nfl_api:nfl_pass@tcp(192.168.2.101:3306)/nfl?autocommit=true")
    //ctx = db.OpenSQL(ctx, "nfl", "mysql", "root:root@tcp(172.16.102.129:3306)/nfl?autocommit=true")
    defer db.Close(ctx) // closes all DB connections
    var wg sync.WaitGroup
    wg.Add(1)
    seattle := nfl.GetTeam("SEA", ctx, wg)
    jsonStr, _ := json.MarshalIndent(seattle, "", "\t")
    os.Stdout.Write(jsonStr)

    start := time.Now()


    iter := 2

    wg.Add(iter)

    loop := make([]nfl.Teams, iter, iter)

    for i := 0; i < iter; i++ {
        loop[i] = nfl.GetTeams("", ctx, wg)

    }

    elapsed := time.Since(start)
    fmt.Printf("\n%v database proc executions took %s\n",iter, elapsed)

    teamTotal := 0
    for i := 0; i < len(loop); i++ {
        teamTotal += len(loop[i])
    }
    fmt.Printf("rows retreived: %v", teamTotal)
    //jsonStr, _ := json.MarshalIndent(loop, "", "\t")
    //os.Stdout.Write(jsonStr)



}

func trackExecutionTime(start time.Time, name string) {
    elapsed := time.Since(start)
    log.Printf("%s took %s", name, elapsed)
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
