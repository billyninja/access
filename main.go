package main

/*
    filtering engine:
       methods, url prefixes, url suffixes

    Detect Time Window
        1 - parse datetime
*/


import (
    "fmt"
    "time"
    "io/ioutil"
    "strings"
    "sort"
)

var (
    filter *Filters
)

const (
    DATE_FMT_A = "02/Jan/2006:15:04:05 -0700"
)

type Filters struct {
    Methods []string
    Prefixes []string
    Sufixes []string
}

func (f *Filters) Filter(le *LogEntry) bool {
    var valid bool
    for _, mt := range f.Methods {
        if le.Method == mt {
            valid = true
            break
        }
    }

    for _, pf := range f.Prefixes {
        if strings.HasPrefix(le.Url, pf) {
            valid = true
            break
        }
    }

    for _, sf := range f.Sufixes {
        if strings.HasSuffix(le.Url, sf) {
            valid = true
            break
        }
    }

    return valid
}


type LogEntry struct{
    Time        time.Time
    IP          string
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
    Hits        int
    First       time.Time
    Last        time.Time
    City        string
    Country     string
    PerHour     []*PerHour
    PerUrl      []*UrlCount
}

type ByHits []*Client

func (a ByHits) Len() int {
    return len(a)
}
func (a ByHits) Swap(i, j int)      {
    a[i], a[j] = a[j], a[i]
}
func (a ByHits) Less(i, j int) bool {
    return a[i].Hits < a[j].Hits
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
    fragsE := strings.Split(fragsB[1], "]")
    fragsF := strings.Split(fragsE[0], "[")

    if len(fragsA) < 5 || len(fragsD) < 2 {
        fmt.Printf(">\n\n\n%+v\n\n\n", fragsA)
        return nil
    }


    tm, err := time.Parse(DATE_FMT_A, fragsF[1])
    if err != nil {
        fmt.Printf("Err parsing time:\n%+v\n", err)
        return nil
    }

    le := &LogEntry{
        Time: tm,
        IP: fragsB[0],
        UA: fragsA[5],
        Status: fragsC[0],
        Method: fragsD[0],
        Url: fragsD[1],
        Ref: fragsA[3],
    }
    if filter != nil && !filter.Filter(le) {
        return nil
    }

    return le
}


func main() {

    filter = &Filters{
        Methods: []string{"POST"},
        Prefixes: []string{"/site/aereo"},
        Sufixes: []string{"/1/0/1"},
    }

    CliMap := map[string]*Client{}
    t1 := time.Now()
    content, err := ioutil.ReadFile("sample.txt")
    if err != nil {
        fmt.Printf("Err: %v", err)
        return
    }

    sorted := []*Client{}
    lines := strings.Split(string(content), "\n")

    var first, last *LogEntry
    var count uint64
    for _, ll := range lines {
        le := get_log_entry(ll)
        if le == nil {
            continue
        }
        last = le
        if first == nil {
            first = le
        }
        count++

        mkey := fmt.Sprintf("%s---%s", le.IP, le.UA)
        if val, ok := CliMap[mkey]; ok {
            val.Hits += 1
            val.Last = le.Time
        } else {
            cli := &Client{
                IP:  le.IP,
                UA: le.UA,
                Hits: 1,
                First: le.Time,
            }
            CliMap[mkey] = cli
            sorted = append(sorted, cli)
        }
    }
    sort.Sort(sort.Reverse(ByHits(sorted)))
    fmt.Printf("\n%s\n", time.Since(t1))

    fmt.Printf("Captured %d hits between %s and %s, counted %d (%.2f%%)\n",
               len(lines), first.Time.Format("01/02 15:04"), last.Time.Format("01/02 15:04"),
               count, (float64(count)/float64(len(lines))*100.0))

    fmt.Printf("==== \nTOP - 10 HITTERS \n")
    for i := 0; i < 10; i++ {
        fmt.Printf("#%d - %d - %s\n%s\n\n", (i + 1), sorted[i].Hits, sorted[i].IP, sorted[i].UA)
    }
}
