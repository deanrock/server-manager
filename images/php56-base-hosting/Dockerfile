FROM deanrock/server-manager:php56-base

ADD php-fpm.conf /etc/php5/fpm/php-fpm.conf
ADD www.conf /etc/php5/fpm/pool.d/www.conf
ADD php.ini /etc/php5/fpm/php.ini

RUN mkdir /mystuff

ADD prepend.php /mystuff/prepend.php

RUN mv /usr/sbin/php5-fpm /usr/bin/php5-fpm
