package main

import (
    "io"
    "gopkg.in/xmlpath.v2"
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