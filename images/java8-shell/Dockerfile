FROM manager/java8-base

RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd

RUN mkdir /mystuff
ADD start.sh /mystuff/

USER #user#

ENV PATH $PATH:/usr/local/apache-maven/apache-maven-3.2.5/bin
ENV JAVA_HOME /usr/lib/jvm/java-8-oracle
