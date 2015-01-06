#!/bin/bash

sed "s/PROFILER_PORT_9000_TCP_ADDR/$PROFILER_PORT_9000_TCP_ADDR/" www.conf > /etc/nginx/sites-enabled/www.conf

nginx
