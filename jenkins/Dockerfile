FROM ubuntu:16.04

RUN apt-get update
RUN apt-get install -y curl software-properties-common git apt-transport-https

##
## Go
##
RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update
RUN apt-get install -y golang-go
ENV GOPATH=$HOME/go
ENV PATH=$PATH:$HOME/go/bin

##
## Go tools
##
RUN go get -u github.com/kardianos/govendor
RUN go get -u github.com/jstemmer/go-junit-report
RUN go get -u github.com/t-yuki/gocover-cobertura

##
## MongoDB
##
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 9DA31620334BD75D9DCB49F368818C72E52529D4
RUN echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu xenial/mongodb-org/4.0 multiverse" \
    | tee /etc/apt/sources.list.d/mongodb-org-4.0.list
RUN apt-get update
RUN apt-get install -y mongodb-org

COPY entry.sh /entry.sh
RUN chmod +x /entry.sh
ENTRYPOINT ["/entry.sh"]

