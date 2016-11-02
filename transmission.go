package main

import (
    "sync"
    "net/http"
    "encoding/json"
    "bytes"
    "io/ioutil"
)

const rpcUrl = "http://localhost/transmission/rpc"
const sessionHdr = "X-Transmission-Session-Id"

var (
    wg sync.WaitGroup
    sessionId string
    trBasicAuth string
)

type TrRequestArgs struct {
    Filename string  `json:"filename"`
}

type TrRequest struct {
    Method string `json:"method"`
    Arguments TrRequestArgs `json:"arguments"`
}

func spawnTransmissionSession(auth string) {
    wg.Add(1)
    go func() {
        defer wg.Done()
        getTransmissionSession(auth)
    }()
}

func getTransmissionSession(auth string) {
    trBasicAuth = auth

    req, err := http.NewRequest("GET", rpcUrl, nil)
    if err != nil { panic(err) }
    setBasicAuth(req, trBasicAuth)

    resp, err := (&http.Client{}).Do(req)
    defer resp.Body.Close()
    if err != nil { panic(err) }

    sessionId = resp.Header.Get(sessionHdr)
}

func addToTransmission(magnet string) (string, error) {
    wg.Wait()

    data := TrRequest{
        Method: "torrent-add",
        Arguments: TrRequestArgs{
            Filename: magnet,
        },
    }
    jsonData, err := json.Marshal(data)
    if err != nil { return "", err }

    req, err := http.NewRequest("POST", rpcUrl, bytes.NewBufferString(string(jsonData)))
    if err != nil { return "", err }

    setBasicAuth(req, trBasicAuth)
    req.Header.Add(sessionHdr, sessionId)

    client := &http.Client{}
    resp, err := client.Do(req)
    defer resp.Body.Close()
    if err != nil { return "", err }

    body, err := ioutil.ReadAll(resp.Body)

    return string(body), err
}
