package main

import (
    "net/http"
    "os"
    "io/ioutil"
    "strings"
)

type Episode struct {
    title string
    magnetUri string
}

func airedToday(feedLink, seenFile string) ([]Episode, error) {
    var episodes []Episode
    var seen []string

    resp, err := http.Get(feedLink)
    if err != nil { return episodes, err }
    defer resp.Body.Close()

    seenContent, err := ioutil.ReadFile(seenFile)
    if err == nil {
        seen = strings.Split(string(seenContent), "\n")
    } else if !os.IsNotExist(err) {
        return episodes, err
    }

    items, err := xpath(resp.Body, "//item")
    if err != nil { return episodes, err }
    for _, item := range items {
        if !existsInList(seen, xpathN(item, ".//guid")) {
            episodes = append(episodes, Episode{title: xpathN(item, ".//title"), magnetUri: xpathN(item, ".//link")})
            seen = append(seen, xpathN(item, ".//guid"))
        }
    }

    err = ioutil.WriteFile(seenFile, []byte(strings.Join(seen, "\n")), 0644)
    if err != nil { return episodes, err }

    return episodes, nil
}
