package main

import (
    "fmt"
    "os"
    "sync"
    "regexp"
)

const magnetPrefix = "magnet:?xt=urn:btih:"

func download(title string) {
    titleRegex := regexp.MustCompile(`.* S\d\dE\d\d`)
    replaceRegex := regexp.MustCompile(`\s*\(.*\)`)
    title = replaceRegex.ReplaceAllLiteralString(titleRegex.FindString(title), "")

    fmt.Print("Searching '" + title + "'... ")

    hash, err := searchTorrent(title)
    if err != nil { panic(err) }

    fmt.Println(hash)

    magnet := magnetPrefix + hash
    resp, err := addToTransmission(magnet)
    if err != nil { panic(err) }

    fmt.Println(resp)
}

func main() {
    spawnTransmissionSession()

    episodes, err := airedToday(os.Args[1])
    if err != nil { panic(err) }

    var wg sync.WaitGroup
    wg.Add(len(episodes))

    for _, title := range episodes {
        go func(title string) {
            defer wg.Done()
            download(title)
        }(title)
    }

    wg.Wait()
}
