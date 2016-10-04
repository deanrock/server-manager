FROM debian:wheezy

RUN apt-get clean && apt-get update && apt-get install -y ssmtp  build-essential libncurses5-dev libncursesw5-dev libreadline6-dev libdb5.1-dev libgdbm-dev libsqlite3-dev libssl-dev libbz2-dev libexpat1-dev liblzma-dev zlib1g-dev wget ca-certificates git subversion mercurial git subversion mercurial strace tar wget curl gcc build-essential libmysqlclient-dev mysql-client screen nano joe vim unzip procps fontconfig libxrender1 xfonts-base xfonts-75dpi libpng12-dev libjpeg-dev openssh-server libopenjpeg-dev libtiff-dev libfreetype6-dev libwebp-dev telnet netcat libxml2-dev libxslt1-dev libcairo2 libpango1.0-0 libgdk-pixbuf2.0-0 libffi-dev shared-mime-info zip libmagickwand-dev

RUN echo "Europe/Ljubljana" > /etc/timezone; dpkg-reconfigure -f noninteractive tzdata

ADD ssmtp.conf /etc/ssmtp/ssmtp.conf

#ffmpeg
RUN echo "deb http://www.deb-multimedia.org wheezy main non-free" >> /etc/apt/sources.list && \
	echo "deb-src http://www.deb-multimedia.org wheezy main non-free" >> /etc/apt/sources.list && \
	echo "deb http://apt.postgresql.org/pub/repos/apt/ wheezy-pgdg main" >> /etc/apt/sources.list && \
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

RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | \
  apt-key add -
RUN apt-get update

#wkhtmltox
RUN apt-get install -y fontconfig libxrender1 xfonts-base xfonts-75dpi && \
	cd /tmp && \
	wget -O wkhtmltox.deb  http://download.gna.org/wkhtmltopdf/0.12/0.12.2.1/wkhtmltox-0.12.2.1_linux-wheezy-amd64.deb && \
	dpkg -i wkhtmltox.deb

RUN mkdir -p /usr/local/opt && mkdir -p /usr/local/src && wget -O /usr/local/src/python344.tgz https://www.python.org/ftp/python/3.4.4/Python-3.4.4.tgz \
    && cd /usr/local/src && tar xvfz python344.tgz && cd /usr/local/src/Python-3.4.4 && ./configure --prefix=/usr/local/opt/python-3.4 \
    && make && make install

RUN apt-get install -y libffi-dev libmysqlclient-dev libmemcached-dev libpq-dev python-psycopg2 postgresql-client-9.5

RUN mkdir -p /usr/local/opt && mkdir -p /usr/local/src && wget -O /usr/local/src/python351.tgz https://www.python.org/ftp/python/3.5.1/Python-3.5.1.tgz \
    && cd /usr/local/src && tar xvfz python351.tgz && cd /usr/local/src/Python-3.5.1 && ./configure --prefix=/usr/local/opt/python-3.5 \
    && make && make install
    
ENV LANG "C.utf8"
ENV LC_ALL "C.UTF-8"
ENV LC_LANG "C.UTF-8"
