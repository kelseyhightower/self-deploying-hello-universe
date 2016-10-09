FROM alpine:3.1
MAINTAINER Kelsey Hightower <kelsey.hightower@gmail.com>
RUN mkdir -p /opt/bin
ADD hello-universe /opt/bin/hello-universe
