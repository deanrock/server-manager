FROM deanrock/server-manager:python27-base

RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd

RUN mkdir /mystuff

ADD start.sh /mystuff/start.sh

RUN apt-get install locales && locale-gen "en_US.UTF-8" && dpkg-reconfigure locales

RUN chown -R #user#:#user# /mystuff && chmod +x /mystuff/start.sh

USER #user#

EXPOSE 9000

CMD /bin/bash -c "/mystuff/start.sh"
