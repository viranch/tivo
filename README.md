# tivo
Auto download recently aired episodes of TV shows you care about

This is basically a Go version of [nning/transmission-rss](https://github.com/nning/transmission-rss).

## Run
```
go get github.com/viranch/tivo gopkg.in/xmlpath.v2
go install github.com/viranch/tivo
$GOPATH/bin/tivo
```

This program is intended to be used with/inside [docker-tv](https://github.com/viranch/docker-tv):

```
docker run -d --name tv -p 80:80 viranch/tv
$GOPATH/bin/tivo -feed 'http://showrss.info/user/XXXXXX.rss?magnets=true'
```

Then browse to [http://localhost](http://localhost) to see today's aired episodes added to downloads.
