#!/bin/bash -eu
CWD=$(cd `dirname $0`/../;pwd)
TARGET=${1:-"-func"}

cd $CWD
[ -e $CWD/.coverage ] || mkdir $CWD/.coverage
FILE=".coverage/coverage.out"
go test ./... -cover -coverprofile=${FILE}
go tool cover -func=${FILE}
