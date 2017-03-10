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

func download(title string, trWg *sync.WaitGroup) error {
    titleRegex := regexp.MustCompile(`.* S\d\dE\d\d`)
    replaceRegex := regexp.MustCompile(`\s*\(.*\)`)
    title = replaceRegex.ReplaceAllLiteralString(titleRegex.FindString(title), "")

    if searchSuffix != "" {
        title = title + " " + searchSuffix
    }

    fmt.Println(title)

    hash, err := searchTorrent(title, basicAuth)
    if err != nil || hash == "" {
        // fallback to skytorrents.in
        hash, err = searchSkyTorrents(title, basicAuth)
    }
    if err != nil { return err }
    if hash == "" { return fmt.Errorf("No torrent found") }

    fmt.Println(title, ":", hash)

    trWg.Wait()
    magnet := magnetPrefix + hash
    resp, err := addToTransmission(magnet)
    if err != nil { return err }

    fmt.Println(resp)

    return nil
}

func fatal(err error) {
    panic(err)
}

func main() {
    initFlags()

    if episodeFeedLink == "" {
        fmt.Println("No episode RSS feed specified\n")
        flag.Usage()
        os.Exit(1)
    }

    var trWg sync.WaitGroup
    trWg.Add(1)
    go func(auth string) {
        defer trWg.Done()
        err := getTransmissionSession(auth)
        if err != nil { fatal(err) }
    }(basicAuth)

    episodes, err := airedToday(episodeFeedLink)
    if err != nil { fatal(err) }

    var wg sync.WaitGroup
    wg.Add(len(episodes))

    for _, title := range episodes {
        go func(title string) {
            defer wg.Done()
            err := download(title, &trWg)
            if err != nil {
                fmt.Printf("Error processing '%s': %s\n", title, err)
            }
        }(title)
    }

    wg.Wait()
}
