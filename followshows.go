package main

import (
    "net/http"
    "time"
)

func airedToday(feedLink string) ([]string, error) {
    var episodes []string

    resp, err := http.Get(feedLink)
    if err != nil { return episodes, err }
    defer resp.Body.Close()

    today := time.Now().Format("02 Jan 2006")

    return xpathS(resp.Body, "//item/pubDate[contains(text(), '" + today + "')]/../title/text()")
}
