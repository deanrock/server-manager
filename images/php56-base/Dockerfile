FROM manager/debian7basehosting

RUN echo "deb http://packages.dotdeb.org wheezy-php56 all" >> /etc/apt/sources.list.d/dotdeb.list && \
	echo "deb-src http://packages.dotdeb.org wheezy-php56 all" >> /etc/apt/sources.list.d/dotdeb.list && \
	apt-get update && apt-get install -y wget && \
	wget http://www.dotdeb.org/dotdeb.gpg -O- |apt-key add - && apt-get update && apt-get install -y php5-cli php5-fpm \
	php5-curl php5-gd php5-intl php-pear php5-imagick php5-imap php5-mcrypt php5-memcache php5-pspell php5-recode php5-sqlite php5-tidy php5-xmlrpc php5-xsl \
	php5-dev

RUN apt-get install -y php5-gd php5-mysql

#RUN pear clear-cache && pear update-channels
RUN pecl install xhprof-beta

RUN curl -sS https://getcomposer.org/installer | php && mv composer.phar /usr/local/bin/composer

RUN pear channel-discover pear.drush.org && pear install drush/drush

RUN drush
