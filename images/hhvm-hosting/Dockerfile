FROM manager/debian7basehosting

RUN echo deb http://dl.hhvm.com/debian wheezy main | tee /etc/apt/sources.list.d/hhvm.list && \
	apt-get update && apt-get install -y --force-yes hhvm

EXPOSE 9000

RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd
USER #user#

CMD hhvm --mode server -vServer.Type=fastcgi -vServer.Port=9000
