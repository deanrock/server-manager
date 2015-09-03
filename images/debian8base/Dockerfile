FROM debian:jessie

RUN apt-get update && apt-get install -y ssmtp  build-essential libncurses5-dev libncursesw5-dev libreadline6-dev libgdbm-dev libsqlite3-dev libssl-dev libbz2-dev libexpat1-dev liblzma-dev zlib1g-dev wget ca-certificates git subversion mercurial git subversion mercurial strace tar wget curl gcc build-essential libmysqlclient-dev mysql-client screen nano joe vim unzip procps fontconfig libxrender1 xfonts-base xfonts-75dpi libpng12-dev libjpeg-dev openssh-server libopenjpeg-dev libtiff-dev libfreetype6-dev libwebp-dev telnet netcat libxml2-dev libxslt1-dev libcairo2 libpango1.0-0 libgdk-pixbuf2.0-0 libffi-dev shared-mime-info zip

RUN echo "Europe/Ljubljana" > /etc/timezone; dpkg-reconfigure -f noninteractive tzdata

