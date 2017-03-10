package main

import (
    "net/http"
    "regexp"
    "strconv"
    "net/url"
    "io/ioutil"
)

func searchSkyTorrents(title, auth string) (string, error) {
    req, err := http.NewRequest("GET", "http://localhost/st/rss/all/ad/1/" + url.PathEscape(title), nil)
    if err != nil { return "", err }

    setBasicAuth(req, auth)
    pretendToBeChrome(req)

    resp, err := (&http.Client{}).Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil { return "", err }

    torrentsDesc, err := xpathR(body, "//item/description/text()")
    if err != nil { return "", err }
    torrentsHash, err := xpathR(body, "//item/guid/text()")
    if err != nil { return "", err }

    descRegex, err := regexp.Compile(`(\d+) seeder\(s\), (\d+) leecher\(s\), `)
    if err != nil { return "", err }
    hashRegex, err := regexp.Compile(`/info/(\w+)/`)
    if err != nil { return "", err }

    score := 0
    var winner string
    for i, desc := range torrentsDesc {
        matches := descRegex.FindStringSubmatch(desc)

        seeds, err := strconv.Atoi(matches[1])
        if err != nil { return "", err }
        peers, err := strconv.Atoi(matches[2])
        if err != nil { return "", err }

        newScore := (seeds * 2) + peers
        if newScore > score {
            winner = hashRegex.FindStringSubmatch(torrentsHash[i])[1]
            score = newScore
        }
    }

    return winner, nil
}
