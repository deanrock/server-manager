FROM deanrock/server-manager:debian8base

RUN echo "deb http://repos.zend.com/zend-server/early-access/php7/repos ubuntu/" >> /etc/apt/sources.list && \
  apt-get update && \
  apt-get install -y --force-yes php7-beta1
