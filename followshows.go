package main

import (
    "net/http"
    "time"
)

func airedToday(feedLink string) ([]string, error) {
    var episodes []string

    resp, err := http.Get(feedLink)
    defer resp.Body.Close()
    if err != nil { return episodes, err }

    today := time.Now().Format("02 Jan 2006")

    return xpath(resp.Body, "//item/pubDate[contains(text(), '" + today + "')]/../title/text()")
}
