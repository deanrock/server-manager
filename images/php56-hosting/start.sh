mkdir -p /home/#user#/apps/#appname#/logs && chown -R #user#:#user# /home/#user#/apps/#appname#/logs

echo "error_log = /home/#user#/apps/#appname#/logs/php5-fpm.log" >> /etc/php5/fpm/pool.d/www.conf

php5-fpm
