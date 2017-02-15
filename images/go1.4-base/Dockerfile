FROM deanrock/server-manager:debian7base

# gcc for cgo
RUN apt-get update && apt-get install -y \
        gcc libc6-dev make \
        --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.4.1

RUN curl -sSL https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz \
        | tar -v -C /usr/src -xz

RUN cd /usr/src/go/src && ./make.bash --no-clean 2>&1

ENV PATH /usr/src/go/bin:$PATH
