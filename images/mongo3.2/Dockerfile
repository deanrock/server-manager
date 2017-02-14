FROM deanrock/server-manager:mongo3.2-base

RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd

RUN mkdir /mystuff

ADD start.sh /mystuff/start.sh
ADD mongod.conf /etc/mongod.conf

RUN chown -R #user#:#user# /mystuff && chmod +x /mystuff/start.sh

USER #user#

EXPOSE 27017

CMD [ "/mystuff/start.sh" ]
