FROM deanrock/server-manager:go1.4-base

RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd

RUN mkdir /mystuff

ADD start.sh /mystuff/start.sh

RUN chown -R #user#:#user# /mystuff && chmod +x /mystuff/start.sh

USER #user#

EXPOSE 8080

CMD /bin/bash -c "/mystuff/start.sh"
