package main

import (
    "net/http"
    "regexp"
    "strconv"
)

func searchTorrent(title, auth string) (string, error) {
    req, err := http.NewRequest("GET", "http://localhost/tz/feed", nil)
    if err != nil { return "", err }

    q := req.URL.Query()
    q.Add("f", title)
    req.URL.RawQuery = q.Encode()
    setBasicAuth(req, auth)

    client := &http.Client{}
    resp, err := client.Do(req)
    defer resp.Body.Close()
    if err != nil { return "", err }

    torrents, err := xpath(resp.Body, "//item/description/text()")
    if err != nil { return "", err }

    descRegex, err := regexp.Compile(`Size: (\d+ \w+) Seeds: (\d+) Peers: (\d+) Hash: (\w+)`)
    if err != nil { return "", err }

    score := 0
    var winner string
    for _, torrent := range torrents {
        matches := descRegex.FindStringSubmatch(torrent)

        seeds, err := strconv.Atoi(matches[2])
        if err != nil { return "", err }
        peers, err := strconv.Atoi(matches[3])
        if err != nil { return "", err }

        if (seeds * 2) + peers > score {
            winner = matches[4]
        }
    }

    return winner, nil
}
