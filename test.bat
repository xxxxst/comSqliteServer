echo off

set GOARCH=386
set CGO_ENABLED=1

go run "src/test/ComTest.go"
