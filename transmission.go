package main

import (
    "net/http"
    "encoding/json"
    "bytes"
    "io/ioutil"
)

const rpcUrl = "/transmission/rpc"
const sessionHdr = "X-Transmission-Session-Id"

var (
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

func getTransmissionSession(remote, auth string) error {
    trBasicAuth = auth

    req, err := http.NewRequest("GET", remote + rpcUrl, nil)
    if err != nil { return err }
    setBasicAuth(req, trBasicAuth)

    resp, err := (&http.Client{}).Do(req)
    if err != nil { return err }
    defer resp.Body.Close()

    sessionId = resp.Header.Get(sessionHdr)

    return nil
}

func addToTransmission(remote, magnet string) (string, error) {
    data := TrRequest{
        Method: "torrent-add",
        Arguments: TrRequestArgs{
            Filename: magnet,
        },
    }
    jsonData, err := json.Marshal(data)
    if err != nil { return "", err }

    req, err := http.NewRequest("POST", remote + rpcUrl, bytes.NewReader(jsonData))
    if err != nil { return "", err }

    setBasicAuth(req, trBasicAuth)
    req.Header.Add(sessionHdr, sessionId)

    resp, err := (&http.Client{}).Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)

    return string(body), err
}
