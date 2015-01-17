FROM manager/debian7basehosting

RUN mkdir -p /usr/local/opt && mkdir -p /usr/local/src && wget -O /usr/local/src/python279.tgz https://www.python.org/ftp/python/2.7.9/Python-2.7.9.tgz \
    && cd /usr/local/src && tar xvfz python279.tgz && cd /usr/local/src/Python-2.7.9 && ./configure --prefix=/usr/local/opt/python-2.7 \
    && make && make install
