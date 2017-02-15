FROM deanrock/server-manager:memcached-base

EXPOSE 11211

CMD ["-m", "64"]

#user stuff
RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd

USER #user#

ENTRYPOINT memcached
