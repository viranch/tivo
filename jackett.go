package main

import (
    "net/http"
    "strconv"
    "net/url"
    "io/ioutil"
    "strings"
)

func searchJackett(title, auth, apiKey string) (string, error) {
    url := "http://localhost/jk/api/v2.0/indexers/all/results/torznab/api?t=search&q=" + url.PathEscape(title) + "&apikey=" + url.PathEscape(apiKey)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil { return "", err }

    setBasicAuth(req, auth)
    pretendToBeChrome(req)

    resp, err := (&http.Client{}).Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil { return "", err }

    bodyStr := string(body[:])
    replacedBody := strings.Replace(bodyStr, "torznab:attr", "torznabAttr", -1)
    body = []byte(replacedBody)
    items, err := xpathByte(body, "//item")
    if err != nil { return "", err }

    score := 0
    var winner string
    for _, item := range items {
        seeds, err := strconv.Atoi(xpathN(item, ".//torznabAttr[@name='seeders']/@value"))
        if err != nil { continue }
        peers, err := strconv.Atoi(xpathN(item, ".//torznabAttr[@name='peers']/@value"))
        if err != nil { continue }

        newScore := (seeds * 2) + peers
        if newScore > score {
            winner = xpathN(item, ".//torznabAttr[@name='magneturl']/@value")
            score = newScore
        }
    }

    return winner, nil
}
