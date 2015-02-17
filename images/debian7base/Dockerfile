FROM debian:wheezy

RUN apt-get update && apt-get install -y ssmtp  build-essential libncurses5-dev libncursesw5-dev libreadline6-dev libdb5.1-dev libgdbm-dev libsqlite3-dev libssl-dev libbz2-dev libexpat1-dev liblzma-dev zlib1g-dev wget ca-certificates git subversion mercurial git subversion mercurial strace tar wget curl gcc build-essential libmysqlclient-dev mysql-client screen nano joe vim unzip procps

RUN echo "Europe/Ljubljana" > /etc/timezone; dpkg-reconfigure -f noninteractive tzdata

ADD ssmtp.conf /etc/ssmtp/ssmtp.conf
