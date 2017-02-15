FROM deanrock/server-manager:redis-base

#start script
RUN mkdir /mystuff
ADD start.sh /mystuff/start.sh

#user stuff
RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd
RUN chown -R #user#:#user# /mystuff && chmod +x /mystuff/start.sh

USER #user#

ADD redis.conf /mystuff/redis.conf

EXPOSE 6379
CMD [ "/mystuff/start.sh" ]
