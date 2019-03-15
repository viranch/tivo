package main

import (
    "fmt"
    "os"
    "flag"
    "sync"
)

var (
    basicAuth string
    episodeFeedLink string
    remote string
    seenFile string
)

func initFlags() {
    flag.StringVar(&episodeFeedLink, "feed", "", "Link to episode RSS feed")
    flag.StringVar(&remote, "remote", "http://127.0.0.1", "Transmission and search remote base URL, default: http://127.0.0.1")
    flag.StringVar(&basicAuth, "auth", "", "Basic authentication credentials in USER:PASSWORD format")
    flag.StringVar(&seenFile, "seen", os.Getenv("HOME") + "/.config/tivo.seen", "Location of the seen file, default: ~/.config/tivo.seen")
    flag.Parse()
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
        err := getTransmissionSession(remote, auth)
        if err != nil { fatal(err) }
    }(basicAuth)

    episodes, err := airedToday(episodeFeedLink, seenFile)
    if err != nil { fatal(err) }

    if len(episodes) == 0 {
        fmt.Println("No new episodes\n")
    }

    trWg.Wait()

    for _, episode := range episodes {
        fmt.Println(episode.title, ":", episode.magnetUri)
        resp, err := addToTransmission(remote, episode.magnetUri)
        if err != nil {
            fmt.Println(err)
        } else {
            fmt.Println(resp)
        }
    }
}
