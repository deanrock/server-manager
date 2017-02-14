FROM deanrock/server-manager:debian7base

RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd

RUN mkdir /mystuff

RUN apt-get install locales && locale-gen "en_US.UTF-8" && dpkg-reconfigure locales

ADD start.sh /mystuff/start.sh

RUN chown -R #user#:#user# /mystuff && chmod +x /mystuff/start.sh

USER #user#

EXPOSE 9000

CMD /bin/bash -c "/mystuff/start.sh"
