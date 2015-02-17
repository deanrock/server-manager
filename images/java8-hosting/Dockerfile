FROM manager/java8-base

RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd

RUN mkdir /mystuff

ADD start.sh /mystuff/start.sh

RUN chown -R #user#:#user# /mystuff && chmod +x /mystuff/start.sh

USER #user#

EXPOSE 8080

RUN export PATH=$PATH:/usr/local/apache-maven/apache-maven-3.2.5/bin

CMD /bin/bash -c "/mystuff/start.sh"
