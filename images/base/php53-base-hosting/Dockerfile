FROM ubuntu:12.04

RUN apt-get update && apt-get install -y ssmtp  build-essential libncurses5-dev libncursesw5-dev libreadline6-dev libgdbm-dev libsqlite3-dev libssl-dev libbz2-dev libexpat1-dev liblzma-dev zlib1g-dev wget ca-certificates


RUN apt-get install --force-yes -y wget apt-get install -y php5-cli php5-fpm \
	php5-curl php5-gd php5-intl php-pear php5-imagick php5-imap php5-mcrypt php5-memcache php5-pspell php5-recode php5-sqlite php5-tidy php5-xmlrpc php5-xsl

RUN apt-get install -y php5-gd php5-mysql

ADD ssmtp.conf /etc/ssmtp/ssmtp.conf

ADD php-fpm.conf /etc/php5/fpm/php-fpm.conf
ADD www.conf /etc/php5/fpm/pool.d/www.conf
ADD php.ini /etc/php5/fpm/php.ini
