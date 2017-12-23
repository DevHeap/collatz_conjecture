#!/bin/bash
DIR=$1
echo Working directory: $DIR

HTTPPID=0

function restart() {
    if [[ $HTTPPID != 0 ]]; then
    	kill $HTTPPID
    fi
    pushd src/frontend
    python -m SimpleHTTPServer 8000 &
    HTTPPID=$!
    popd 

    kill -9 $(/sbin/pidof backend)
    pushd src/backend
    GOPATH=$HOME/go go build || exit 1
    ./backend --addr 0.0.0.0:8080 &
    popd
}

restart

CURRENT_HASH=$(git rev-parse HEAD)
while true; do
    echo Checking updates
    git pull
    NEW_HASH=$(git rev-parse HEAD)
    echo old hash: $CURRENT_HASH
    echo new hash: $NEW_HASH
    if [[ $CURRENT_HASH != $NEW_HASH ]]; then
        echo New version availible, restarting server
        CURRENT_HASH=$NEW_HASH
        restart
    fi

    sleep 10
done
