FROM deanrock/server-manager:ruby22-base

RUN apt-get install curl && \
	curl -sL https://deb.nodesource.com/setup_6.x | bash - && \
	apt-get install -y nodejs

RUN npm install -g bower grunt-cli webpack gulp

# install mongo cli
ENV GPG_KEYS \
	DFFA3DCF326E302C4787673A01C4E7FAAAB2461C \
	42F3E95A2C4F08279C4960ADD68FA50FEA312927
RUN set -ex \
	&& for key in $GPG_KEYS; do \
		apt-key adv --keyserver ha.pool.sks-keyservers.net --recv-keys "$key"; \
	done

ENV MONGO_MAJOR 3.2
ENV MONGO_VERSION 3.2.6

RUN echo "deb http://repo.mongodb.org/apt/debian wheezy/mongodb-org/$MONGO_MAJOR main" > /etc/apt/sources.list.d/mongodb-org.list

RUN set -x \
	&& apt-get update \
	&& apt-get install -y \
		pwgen \
		mongodb-org=$MONGO_VERSION \
		mongodb-org-shell=$MONGO_VERSION \
		mongodb-org-mongos=$MONGO_VERSION \
		mongodb-org-tools=$MONGO_VERSION \
	&& rm -rf /var/lib/apt/lists/* \
	&& rm -rf /var/lib/mongodb \
	&& mv /etc/mongod.conf /etc/mongod.conf.orig
