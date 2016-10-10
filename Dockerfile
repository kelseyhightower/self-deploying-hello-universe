FROM scratch
MAINTAINER Kelsey Hightower <kelsey.hightower@gmail.com>
ADD hello-universe /hello-universe
ENTRYPOINT ["/hello-universe"]
