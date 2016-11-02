package main

import (
    "fmt"
    "os"
    "flag"
    "sync"
    "regexp"
)

const magnetPrefix = "magnet:?xt=urn:btih:"

var (
    basicAuth string
    episodeFeedLink string
    searchSuffix string
)

func initFlags() {
    flag.StringVar(&basicAuth, "auth", "", "Basic authentication credentials in USER:PASSWORD format")
    flag.StringVar(&episodeFeedLink, "feed", "", "Link to episode RSS feed")
    flag.StringVar(&searchSuffix, "suffix", "", "Torrent search suffix (eg, 720p)")
    flag.Parse()
}

func download(title, suffix, auth string) {
    titleRegex := regexp.MustCompile(`.* S\d\dE\d\d`)
    replaceRegex := regexp.MustCompile(`\s*\(.*\)`)
    title = replaceRegex.ReplaceAllLiteralString(titleRegex.FindString(title), "")

    if suffix != "" {
        title = title + " " + suffix
    }

    fmt.Println(title)

    hash, err := searchTorrent(title, auth)
    if err != nil { panic(err) }

    fmt.Println(title, ":", hash)

    magnet := magnetPrefix + hash
    resp, err := addToTransmission(magnet)
    if err != nil { panic(err) }

    fmt.Println(resp)
}

func main() {
    initFlags()

    if episodeFeedLink == "" {
        fmt.Println("No episode RSS feed specified\n")
        flag.Usage()
        os.Exit(1)
    }

    spawnTransmissionSession(basicAuth)

    episodes, err := airedToday(episodeFeedLink)
    if err != nil { panic(err) }

    var wg sync.WaitGroup
    wg.Add(len(episodes))

    for _, title := range episodes {
        go func(title, suffix, auth string) {
            defer wg.Done()
            download(title, suffix, auth)
        }(title, searchSuffix, basicAuth)
    }

    wg.Wait()
}
