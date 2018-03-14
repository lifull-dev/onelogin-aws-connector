#!/bin/bash -eu
cd $(cd `dirname $0`/../;pwd)

_GOOS=("linux" "darwin" "windows")
_GOARCH=("amd64" "amd64" "amd64")

for((i=0;i<${#_GOOS[@]};++i))
do
  _goos=${_GOOS[$i]}
  _goarch=${_GOARCH[$i]}
  _out=onelogin-aws-connector
  _zip=${_out}_${_goos}_${_goarch}.zip
  GOOS=${_goos} GOARCH=${_goarch} go build -o ${_out} -tags includeClientToken
  zip ${_zip} ${_out}
  rm -f ${_out}
done
