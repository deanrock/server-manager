FROM debian:wheezy

RUN apt-get update && \
    apt-get install -y wget ca-certificates build-essential

RUN wget https://storage.googleapis.com/golang/go1.7.5.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.7.5.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin
