# !/bin/zsh
cd cmd/web
go run $(ls -1 *.go | grep -v _test.go) -addr=":8080"