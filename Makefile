@PHONY: build-linux build-win deploy
build-linux:
	cd src && GOOS=linux GOARCH=amd64 go build -o ../hooks/checkout-bin

build-win:
	cd src && GOOS=windows GOARCH=amd64 go build -o ../hooks/checkout-win.exe

deploy: build-win build-linux