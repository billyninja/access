package main

import (
    "fmt"
    "time"
    "io/ioutil"
    "strings"
    "sort"
)

type LogEntry struct{
    Time        time.Time
    IP    string
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

type ByHits []*Client

func (a ByHits) Len() int           { return len(a) }
func (a ByHits) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHits) Less(i, j int) bool { return a[i].Hits < a[j].Hits }

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
        IP: fragsB[0],
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

    sorted := []*Client{}
    lines := strings.Split(string(content), "\n")
    for _, ll := range lines {
        le := get_log_entry(ll)
        if le == nil {
            continue
        }

        mkey := fmt.Sprintf("%s---%s", le.IP, le.UA)
        if val, ok := CliMap[mkey]; ok {
            val.Hits += 1
        } else {
            cli := &Client{
                IP:  le.IP,
                UA: le.UA,
                Hits: 1,
            }
            CliMap[mkey] = cli
            sorted = append(sorted, cli)
        }
    }
    sort.Reverse(ByHits(sorted))
    fmt.Printf("\n%s\n", time.Since(t1))
    println("=== Top 10 Hitters ===")
    for i := 0; i < 10; i++ {
        println(sorted[1].Hits)
    }
}
