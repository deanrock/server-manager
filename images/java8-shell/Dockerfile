FROM manager/java8-base-shell

RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd

USER #user#

RUN export PATH=$PATH:/usr/local/apache-maven/apache-maven-3.2.5/bin
