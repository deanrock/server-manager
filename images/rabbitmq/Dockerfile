FROM manager/debian8base

#start script
RUN mkdir /mystuff
ADD start.sh /mystuff/start.sh

#user stuff
RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd
RUN chown -R #user#:#user# /mystuff && chmod +x /mystuff/start.sh

#install rabbitmq
RUN \
  wget -qO - https://www.rabbitmq.com/rabbitmq-release-signing-key-public.asc | apt-key add - && \
  echo "deb http://www.rabbitmq.com/debian/ testing main" > /etc/apt/sources.list.d/rabbitmq.list && \
  apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get install -y rabbitmq-server && \
  rm -rf /var/lib/apt/lists/* && \
  rabbitmq-plugins enable rabbitmq_management && \
  echo "[{rabbit, [{loopback_users, []}]}]." > /etc/rabbitmq/rabbitmq.config

USER #user#

ENV RABBITMQ_LOG_BASE /home/#user#/apps/#appname#/log
ENV RABBITMQ_MNESIA_BASE /home/#user#/apps/#appname#/mnesia

EXPOSE 5672
EXPOSE 15672
CMD [ "/mystuff/start.sh" ]
