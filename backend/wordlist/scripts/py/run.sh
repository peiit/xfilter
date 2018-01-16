#!/bin/bash
HOST=127.0.0.1:8006
LOG_FILE=../../log/flask.log

case $1 in
  dev)
    python webapp.py
    ;;
  start)
    gunicorn -w 4 -b $HOST webapp:app > $LOG_FILE 2>&1 &
    echo "gunicorn service has been started!"
    ;;
  stop)
    kill -9 `ps aux | grep gunicorn | awk '{print $2}'`
    echo "gunicorn service has been stopped!"
    ;;
  *)
    echo "Usage: [dev start stop]"
    ;;
esac

