env GOOS=linux GOARCH=amd64  go build -ldflags "-X main.version=0.1 -X main.minversion=`date -u +.%Y%m%d.%H%M%S`"
# env GOOS=darwin GOARCH=amd64 go build
