FROM golang:1.5.3-wheezy

MAINTAINER cghsystems

RUN go get github.com/tools/godep
RUN go get github.com/onsi/ginkgo/ginkgo


RUN apt-get update && apt-get install -y redis-server

ENV VCAP_SERVICES "{ \
  \"p-redis\": [ \
   { \
    \"credentials\": { \
     \"host\": \"127.0.0.1\",\
     \"password\": \"password\", \
     \"port\": 6379 \
    }, \
    \"label\": \"p-redis\", \
    \"name\": \"godata\", \
    \"plan\": \"dedicated-vm\", \
    \"tags\": [ \
     \"pivotal\", \
     \"redis\" \
    ] \
   } \
  ] \
}"
