FROM deanrock/server-manager:hhvm-base

EXPOSE 9000

RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd
USER #user#

CMD hhvm --mode server -vServer.Type=fastcgi -vServer.Port=9000
