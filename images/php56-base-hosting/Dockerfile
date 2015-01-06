FROM manager/debian7basehosting

RUN echo "deb http://packages.dotdeb.org wheezy-php56 all" >> /etc/apt/sources.list.d/dotdeb.list && \
	echo "deb-src http://packages.dotdeb.org wheezy-php56 all" >> /etc/apt/sources.list.d/dotdeb.list && \
	apt-get update && apt-get install -y wget && \
	wget http://www.dotdeb.org/dotdeb.gpg -O- |apt-key add - && apt-get update && apt-get install -y php5-cli php5-fpm \
	php5-curl php5-gd php5-intl php-pear php5-imagick php5-imap php5-mcrypt php5-memcache php5-pspell php5-recode php5-sqlite php5-tidy php5-xmlrpc php5-xsl

RUN apt-get install -y php5-gd php5-mysql

ADD php-fpm.conf /etc/php5/fpm/php-fpm.conf
ADD www.conf /etc/php5/fpm/pool.d/www.conf
ADD php.ini /etc/php5/fpm/php.ini

RUN mv /usr/sbin/php5-fpm /usr/bin/php5-fpm
