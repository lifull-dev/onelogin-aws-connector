#!/bin/bash -eu
CWD=$(cd `dirname $0`/../;pwd)
TARGET=${1:-"-func"}

cd $CWD
[ -e $CWD/.coverage ] || mkdir $CWD/.coverage
find ${CWD} \
    -name '*.go' \
    -not -path "*/vendor/*" \
    -not -name "interface.go" \
    -not -name "main.go" \
    -not -name "env.go" | \
        sed -E "s@$GOPATH/src/@@g" | \
        xargs -I{} dirname {} | \
        sort | \
        uniq | \
        while read DIR; do
            echo "test: ${DIR}"
            FILE=.coverage/$(echo $DIR | sed -E 's@/@_@g').out
            go test $DIR -cover -coverprofile=${FILE}
            if [ -f $FILE ]; then
                echo "coverage:"
                go tool cover ${TARGET} $FILE
            fi
            echo ""
        done
