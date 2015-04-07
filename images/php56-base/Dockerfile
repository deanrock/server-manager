FROM manager/debian7basehosting

RUN echo "deb http://packages.dotdeb.org wheezy-php56 all" >> /etc/apt/sources.list.d/dotdeb.list && \
	echo "deb-src http://packages.dotdeb.org wheezy-php56 all" >> /etc/apt/sources.list.d/dotdeb.list && \
	apt-get update && apt-get install -y wget && \
	wget http://www.dotdeb.org/dotdeb.gpg -O- |apt-key add - && apt-get update && apt-get install -y php5-cli php5-fpm \
	php5-curl php5-gd php5-intl php-pear php5-imagick php5-imap php5-mcrypt php5-memcache php5-pspell php5-recode php5-sqlite php5-tidy php5-xmlrpc php5-xsl \
	php5-dev ncftp

RUN apt-get install -y php5-gd php5-mysql

#wkhtmltox
RUN apt-get install -y fontconfig libxrender1 xfonts-base xfonts-75dpi && \
	cd /tmp && \
	wget -O wkhtmltox.deb  http://downloads.sourceforge.net/project/wkhtmltopdf/0.12.2.1/wkhtmltox-0.12.2.1_linux-wheezy-amd64.deb && \
	dpkg -i wkhtmltox.deb

#RUN pear clear-cache && pear update-channels
RUN pecl install xhprof-beta

RUN curl -sS https://getcomposer.org/installer | php && mv composer.phar /usr/local/bin/composer

RUN pear channel-discover pear.drush.org && pear install drush/drush

RUN drush

#ffmpeg
RUN echo "deb http://www.deb-multimedia.org wheezy main non-free" >> /etc/apt/sources.list && \
	echo "deb-src http://www.deb-multimedia.org wheezy main non-free" >> /etc/apt/sources.list && \
	apt-get update && \
	apt-get install -y --force-yes deb-multimedia-keyring && \
	apt-get update && \
	apt-get install \
 	 -y                  \
	  libfaad-dev         \
	  faad                \
	  faac                \
	  libfaac0            \
	  libfaac-dev         \
	  libmp3lame-dev      \
	  x264                \
	  libx264-dev         \
	  libxvidcore-dev     \
	  build-essential     \
	  checkinstall  && \
	apt-get build-dep -y ffmpeg && \
	cd /tmp && \
	apt-get source ffmpeg && \
	cd /tmp/ffmpeg-* && \
	./configure \
	--enable-gpl \
	--enable-nonfree \
	--enable-libfaac \
	--enable-libgsm \
	--enable-libmp3lame \
	--enable-libtheora \
	--enable-libvorbis \
	--enable-libx264 \
	--enable-libxvid \
	--enable-zlib \
	--enable-postproc \
	--enable-swscale \
	--enable-pthreads \
	--enable-x11grab \
	--enable-libdc1394 \
	--enable-version3 \
	--enable-libopencore-amrnb \
	--enable-libopencore-amrwb && \
	mkdir -p /usr/local/share/ffmpeg  && \
	checkinstall -D --install=no --pkgname=ffmpeg-full --autodoinst=yes -y && \
	dpkg -i ffmpeg-full_*-1_amd64.deb


ADD cli_php.ini /etc/php5/cli/php.ini
