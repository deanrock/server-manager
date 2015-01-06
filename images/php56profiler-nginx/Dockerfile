FROM debian:wheezy

RUN echo "deb http://http.debian.net/debian wheezy-backports main" >> /etc/apt/sources.list
RUN apt-get update && apt-get install -y -t wheezy-backports nginx && rm /etc/nginx/sites-enabled/default
ADD start.sh /start.sh
ADD www.conf /www.conf
ADD nginx.conf /etc/nginx/nginx.conf
CMD /start.sh

EXPOSE 80
