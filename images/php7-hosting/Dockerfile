FROM manager/php7-base

RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd

RUN mkdir /mystuff && \
  chown #user#:#user# /mystuff

ADD php-fpm.conf /usr/local/php7/etc/php-fpm.conf
ADD www.conf /usr/local/php7/etc/php-fpm.d/www.conf
ADD start.sh /mystuff/start.sh

RUN chmod +x /mystuff/start.sh

USER #user#

EXPOSE 9000
CMD /bin/bash -c "/mystuff/start.sh"
