FROM deanrock/server-manager:debian7base

RUN echo deb http://dl.hhvm.com/debian wheezy main | tee /etc/apt/sources.list.d/hhvm.list && \
	apt-get update && apt-get install -y --force-yes hhvm
