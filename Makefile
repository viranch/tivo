default: build

build:
	@GOOS=linux GOARCH=arm go build -o ./dist/tivo && tar zcf ./dist/tivo-linux-armv7-$(ver).tar.gz -C ./dist tivo; \
	 GOOS=linux GOARCH=amd64 go build -o ./dist/tivo && tar zcf ./dist/tivo-linux-amd64-$(ver).tar.gz -C ./dist tivo
