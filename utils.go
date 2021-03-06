package main

import (
    "io"
    "gopkg.in/xmlpath.v2"
    "net/http"
    "encoding/base64"
)

const chromeUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.93 Safari/537.36"

func xpath(r io.Reader, spath string) ([]*xmlpath.Node, error) {
    var results []*xmlpath.Node

    path := xmlpath.MustCompile(spath)
    root, err := xmlpath.Parse(r)
    if err != nil { return results, err }

    iter := path.Iter(root)
    for iter.Next() {
        results = append(results, iter.Node())
    }

    return results, nil
}

func xpathN(node *xmlpath.Node, spath string) string {
    s, _ := xmlpath.MustCompile(spath).String(node)
    return s
}

func setBasicAuth(req *http.Request, auth string) {
    if auth != "" {
        req.Header.Add("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(auth)))
    }
}

func pretendToBeChrome(req *http.Request) {
    req.Header.Add("User-Agent", chromeUserAgent)
}

func existsInList(list []string, element string) bool {
    for _, x := range list {
        if x == element {
            return true
        }
    }
    return false
}
