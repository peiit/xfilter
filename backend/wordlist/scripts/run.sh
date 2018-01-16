#!/bin/bash

APP=word_filter
LOG_FILE=log/access.log
PYTHON_DIR=scripts/py

function build(){
  cd $APP && go build
  cd ..
  mv $APP/$APP bin/
  echo "word_filter has been built!"
}

function start(){
  build
  nohup bin/$APP > $LOG_FILE 2>&1 &
}

function stop(){
  pgrep $APP | xargs kill
}

case $1 in
  build)
    build
    ;;
  start)
    start
    echo "$APP has been started!"
    ;;
  stop)
    stop
    echo "$APP has been stopped!"
    ;;
  dev)
    build
    bin/$APP
    echo "$APP has been running!"
    ;;
  prod)
    case $2 in
      start)
        start
        cd $PYTHON_DIR && ./run.sh start
        ;;
      stop)
        stop
        cd $PYTHON_DIR && ./run.sh stop
        ;;
      *)
        echo "Usage: prod [start stop]"
        ;;
    esac
    ;;
  *)
    echo "Usage: [build dev start stop]"
    echo "Deploy: prod [start stop]"
    ;;
esac