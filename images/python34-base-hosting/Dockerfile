FROM manager/debian7basehosting

RUN mkdir -p /usr/local/opt && mkdir -p /usr/local/src && wget -O /usr/local/src/python342.tgz https://www.python.org/ftp/python/3.4.2/Python-3.4.2.tgz \
    && cd /usr/local/src && tar xvfz python342.tgz && cd /usr/local/src/Python-3.4.2 && ./configure --prefix=/usr/local/opt/python-3.4 \
    && make && make install

RUN apt-get install -y libffi-dev libmysqlclient-dev
