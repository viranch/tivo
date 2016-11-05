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

    errChannel chan error
)

func initFlags() {
    flag.StringVar(&basicAuth, "auth", "", "Basic authentication credentials in USER:PASSWORD format")
    flag.StringVar(&episodeFeedLink, "feed", "", "Link to episode RSS feed")
    flag.StringVar(&searchSuffix, "suffix", "", "Torrent search suffix (eg, 720p)")
    flag.Parse()
}

func search(title string) (string, error) {
    titleRegex := regexp.MustCompile(`.* S\d\dE\d\d`)
    replaceRegex := regexp.MustCompile(`\s*\(.*\)`)
    title = replaceRegex.ReplaceAllLiteralString(titleRegex.FindString(title), "")

    if searchSuffix != "" {
        title = title + " " + searchSuffix
    }

    fmt.Println(title)

    hash, err := searchTorrent(title, basicAuth)
    if err != nil { return "", err }

    fmt.Println(title, ":", hash)

    return hash, nil
}

func download(hash string) error {
    magnet := magnetPrefix + hash
    resp, err := addToTransmission(magnet)
    if err != nil { return err }

    fmt.Println(resp)

    return nil
}

func handleError(err error) {
    errChannel <- err
}

func fingersCrossed(wg *sync.WaitGroup) {
    go func() {
        wg.Wait()
        errChannel <- nil
    }()

    if err := <-errChannel; err != nil {
        panic(err)
    }
}

func main() {
    initFlags()

    if episodeFeedLink == "" {
        fmt.Println("No episode RSS feed specified\n")
        flag.Usage()
        os.Exit(1)
    }

    errChannel = make(chan error, 1)

    var trWg sync.WaitGroup
    trWg.Add(1)
    go func(auth string) {
        defer trWg.Done()
        err := getTransmissionSession(auth)
        if err != nil { handleError(err) }
    }(basicAuth)

    episodes, err := airedToday(episodeFeedLink)
    if err != nil { handleError(err) }

    var wg sync.WaitGroup
    wg.Add(len(episodes))

    for _, title := range episodes {
        go func(title string) {
            defer wg.Done()

            hash, err := search(title)
            if err != nil { handleError(err) }

            trWg.Wait()
            err = download(hash)
            if err != nil { handleError(err) }
        }(title)
    }

    fingersCrossed(&wg)
}
