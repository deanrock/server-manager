FROM deanrock/server-manager:php56-base-hosting

RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd

RUN chown -R #user#:#user# /mystuff

ADD www.conf /etc/php5/fpm/pool.d/www.conf

RUN sed -i s/ENV_USER/#user#/ /etc/php5/fpm/php-fpm.conf && sed -i s/ENV_USER/#user#/ /etc/php5/fpm/pool.d/www.conf

USER #user#
EXPOSE 9000
CMD php5-fpm
