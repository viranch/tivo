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

    trWg sync.WaitGroup
    errChannel chan error
    doneChannel chan bool
)

func initFlags() {
    flag.StringVar(&basicAuth, "auth", "", "Basic authentication credentials in USER:PASSWORD format")
    flag.StringVar(&episodeFeedLink, "feed", "", "Link to episode RSS feed")
    flag.StringVar(&searchSuffix, "suffix", "", "Torrent search suffix (eg, 720p)")
    flag.Parse()
}

func download(title, suffix, auth string) error {
    titleRegex := regexp.MustCompile(`.* S\d\dE\d\d`)
    replaceRegex := regexp.MustCompile(`\s*\(.*\)`)
    title = replaceRegex.ReplaceAllLiteralString(titleRegex.FindString(title), "")

    if suffix != "" {
        title = title + " " + suffix
    }

    fmt.Println(title)

    hash, err := searchTorrent(title, auth)
    if err != nil { return err }

    fmt.Println(title, ":", hash)

    trWg.Wait()
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
        close(doneChannel)
    }()

    select {
    case <-doneChannel:
    case err := <-errChannel:
        if err != nil {
            panic(err)
        }
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
    doneChannel = make(chan bool, 1)

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
        go func(title, suffix, auth string) {
            defer wg.Done()
            err := download(title, suffix, auth)
            if err != nil { handleError(err) }
        }(title, searchSuffix, basicAuth)
    }

    fingersCrossed(&wg)
}
