echo off

set GOARCH=386
set CGO_ENABLED=1

go build -o bin/debug/comSqliteServer.exe
