FROM alpine:latest

MAINTAINER Peter Wilmersdorf <ptdorf@gmail.com>

WORKDIR "/opt"

ADD .docker_build/kmk-api /opt/bin/kmk-api
ADD ./templates /opt/templates
ADD ./static /opt/static

CMD ["/opt/bin/kmk-api"]
