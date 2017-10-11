package main

import (
    "fmt"
    "time"
    "io/ioutil"
    "strings"
)

type LogEntry struct{
    Time        time.Time
    ClientIP    string
    UA          string
    Status      string
    Method      string
    Url         string
    Ref         string
}

type UrlCount struct {
    Url     string
    Hits    string
}

type PerHour struct {
    Hour    time.Time
    Hits    string
}

type Client struct {
    IP          string
    UA          string
    Hits        uint32
    City        string
    Country     string
    PerHour     []*PerHour
    PerUrl      []*UrlCount
}

func get_log_entry(line string) *LogEntry {
    fragsA := strings.Split(line, "\"")
    if len(fragsA) == 0 {
        return nil
    }

    fragsB := strings.Split(fragsA[0], " - - ")
    if len(fragsB) < 2 {
        return nil
    }
    fragsC := strings.Split(fragsA[2], " ")
    fragsD := strings.Split(fragsA[1], " ")

    return &LogEntry{
        ClientIP: fragsB[0],
        UA: fragsA[5],
        Status: fragsC[0],
        Method: fragsD[0],
        Url: fragsD[1],
        Ref: fragsA[3],
    }
}


func main() {
    CliMap := map[string]*Client{}
    t1 := time.Now()
    content, err := ioutil.ReadFile("sample.txt")
    if err != nil {
        println("Err: ", err)
        return
    }
    lines := strings.Split(string(content), "\n")
    all := make([]*LogEntry, len(lines))
    for i, ll := range lines {
        all[i] = get_log_entry(ll)
        if all[i] == nil {
            continue
        }

        mkey := fmt.Sprintf("%s---%s", all[i].ClientIP, all[i].UA)
        if val, ok := CliMap[mkey]; ok {
            println(val)
        }
    }
    fmt.Printf("\n%s\n", time.Since(t1))
}
