FROM manager/debian7base

RUN apt-get install curl && \
	curl -sL https://deb.nodesource.com/setup | bash - && \
	apt-get install -y nodejs

RUN npm install -g bower grunt-cli

