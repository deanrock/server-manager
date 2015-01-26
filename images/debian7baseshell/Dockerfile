FROM manager/debian7base

RUN apt-get update && apt-get install -y git subversion mercurial strace tar wget curl gcc build-essential libmysqlclient-dev mysql-client screen nano joe vim unzip

RUN mkdir /mystuff
ADD start.sh /mystuff/
ADD ssmtp.conf /etc/ssmtp/ssmtp.conf
