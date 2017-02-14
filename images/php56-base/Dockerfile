FROM deanrock/server-manager:debian7basehosting

RUN echo "deb http://packages.dotdeb.org wheezy-php56 all" >> /etc/apt/sources.list.d/dotdeb.list && \
	echo "deb-src http://packages.dotdeb.org wheezy-php56 all" >> /etc/apt/sources.list.d/dotdeb.list && \
	apt-get update && apt-get install -y wget && \
	wget http://www.dotdeb.org/dotdeb.gpg -O- |apt-key add - && apt-get update && apt-get install -y php5-cli php5-fpm \
	php5-curl php5-gd php5-intl php-pear php5-imagick php5-imap php5-mcrypt php5-memcache php5-pspell php5-recode php5-sqlite php5-tidy php5-xmlrpc php5-xsl \
	php5-dev ncftp php5-ldap

RUN apt-get install -y php5-gd php5-mysql

#RUN pear clear-cache && pear update-channels
RUN pecl install xhprof-beta

RUN curl -sS https://getcomposer.org/installer | php && mv composer.phar /usr/local/bin/composer

RUN cd /tmp && wget http://files.drush.org/drush.phar && chmod +x drush.phar && \
	mv drush.phar /usr/local/bin/drush && \
	/usr/local/bin/drush

RUN mkdir /tmp/php5-memcached && \
	cd /tmp/php5-memcached && \
	apt-get install -y php5-dev pkg-config php-pear libmemcached-dev && \
	pecl install igbinary && \
	pecl_memcached_ver="2.2.0" && \
	pecl download memcached-${pecl_memcached_ver} && \
	tar xzvf memcached-${pecl_memcached_ver}.tgz && \
	cd memcached-${pecl_memcached_ver}/ && \
	phpize && \
	./configure --enable-memcached-igbinary --disable-memcached-sasl && \
	make && \
	make install

ADD cli_php.ini /etc/php5/cli/php.ini
