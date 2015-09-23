#!/bin/bash

mkdir -p /home/#user#/apps/#appname#/log
mkdir -p /home/#user#/apps/#appname#/mnesia

ulimit -n 1024
cd /var/lib/rabbitmq
exec /usr/lib/rabbitmq/bin/rabbitmq-server "$@"  > "/home/#user#/apps/#appname#/log/startup_log" 2> "/home/#user#/apps/#appname#/log/startup_err"
