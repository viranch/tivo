package main

import (
    "fmt"
    "os"
    "net/http"
    "gopkg.in/xmlpath.v2"

    "io"
    "time"
    "sync"
    "regexp"
    "strconv"
)

func xpath(r io.Reader, spath string) ([]string, error) {
    var results []string

    path := xmlpath.MustCompile(spath)
    root, err := xmlpath.Parse(r)
    if err != nil { return results, err }

    iter := path.Iter(root)
    for iter.Next() {
        results = append(results, iter.Node().String())
    }

    return results, nil
}

func search(title string) (string, error) {
    req, err := http.NewRequest("GET", "http://localhost/tz/feed", nil)
    if err != nil { return "", err }

    q := req.URL.Query()
    q.Add("f", title)
    req.URL.RawQuery = q.Encode()

    client := &http.Client{}
    resp, err := client.Do(req)
    defer resp.Body.Close()
    if err != nil { return "", err }

    torrents, err := xpath(resp.Body, "//item/description/text()")
    if err != nil { return "", err }

    desc_regex, err := regexp.Compile(`Size: (\d+ \w+) Seeds: (\d+) Peers: (\d+) Hash: (\w+)`)
    if err != nil { return "", err }

    score := 0
    var winner string
    for _, torrent := range torrents {
        matches := desc_regex.FindStringSubmatch(torrent)

        seeds, err := strconv.Atoi(matches[2])
        if err != nil { return "", err }
        peers, err := strconv.Atoi(matches[3])
        if err != nil { return "", err }

        torrent_score := (seeds * 2) + peers
        if torrent_score > score {
            winner = matches[4]
        }
    }

    return winner, nil
}

func download(title string) {
    title_regex := regexp.MustCompile(`.* S\d\dE\d\d`)
    replace_regex := regexp.MustCompile(`\s*\(.*\)`)
    title = replace_regex.ReplaceAllLiteralString(title_regex.FindString(title), "")

    hash, err := search(title)
    if err != nil { panic(err) }

    fmt.Println(title, hash)
}

func main() {
    resp, err := http.Get(os.Args[1])
    defer resp.Body.Close()
    if err != nil { panic(err) }

    today := time.Now().Format("02 Jan 2006")
    aired_today, err := xpath(resp.Body, "//item/pubDate[contains(text(), '" + today + "')]/../title/text()")
    if err != nil { panic(err) }

    var wg sync.WaitGroup
    wg.Add(len(aired_today))

    for _, title := range aired_today {
        go func(title string) {
            defer wg.Done()
            download(title)
        }(title)
    }

    wg.Wait()
}
