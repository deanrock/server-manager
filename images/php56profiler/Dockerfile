FROM php56

RUN apt-get install -y git
RUN cd /mystuff && git clone https://github.com/deanrock/xhprof.io.git

VOLUME ["/mystuff/xhprof.io"]
