#!/bin/bash
DIR=$1
echo Working directory: $DIR

kill -9 $(/sbin/pidof backend)
pushd src/backend
GOPATH=$HOME/go go build || exit 1
./backend &


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
	kill -9 $(/sbin/pidof backend)
        pushd src/backend
        GOPATH=$HOME/go go build
        ./backend &
        popd
    fi

    sleep 10
done
