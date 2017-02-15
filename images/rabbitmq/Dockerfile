FROM deanrock/server-manager:rabbitmq-base

#start script
RUN mkdir /mystuff
ADD start.sh /mystuff/start.sh

#user stuff
RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd
RUN chown -R #user#:#user# /mystuff && chmod +x /mystuff/start.sh

USER #user#

ENV RABBITMQ_LOG_BASE /home/#user#/apps/#appname#/log
ENV RABBITMQ_MNESIA_BASE /home/#user#/apps/#appname#/mnesia

EXPOSE 5672
EXPOSE 15672
CMD [ "/mystuff/start.sh" ]
