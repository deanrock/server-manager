FROM deanrock/server-manager:debian7base

#Oracle Java 8
RUN echo "deb http://ppa.launchpad.net/webupd8team/java/ubuntu trusty main" | tee /etc/apt/sources.list.d/webupd8team-java.list && \
echo "deb-src http://ppa.launchpad.net/webupd8team/java/ubuntu trusty main" | tee -a /etc/apt/sources.list.d/webupd8team-java.list && \
apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys EEA14886 && \
apt-get update && \
echo oracle-java8-installer shared/accepted-oracle-license-v1-1 select true | /usr/bin/debconf-set-selections && \
apt-get install -y oracle-java8-installer

#Maven
RUN cd /tmp && \
wget -O apache-maven-3.2.5-bin.tar.gz http://www.apache.si/maven/maven-3/3.2.5/binaries/apache-maven-3.2.5-bin.tar.gz && \
mkdir -p /usr/local/apache-maven && \
cd /usr/local/apache-maven && \
tar xvfz /tmp/apache-maven-3.2.5-bin.tar.gz && \
export PATH=$PATH:/usr/local/apache-maven/apache-maven-3.2.5/bin

ENV PATH $PATH:/usr/local/apache-maven/apache-maven-3.2.5/bin
ENV JAVA_HOME /usr/lib/jvm/java-8-oracle
