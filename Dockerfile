FROM alpine:latest
MAINTAINER Jose Maria Hidalgo Garcia <jhidalgo3@gmail.com>

RUN apk -U add openssl

ENV VERSION v0..0
ENV DOWNLOAD_URL https://github.com/jhidalgo3/dante-cli/releases/download/$VERSION/dante-cli-alpine-linux-amd64-$VERSION.tar.gz

RUN wget -qO- $DOWNLOAD_URL | tar xvz -C /usr/local/bin

ENTRYPOINT ["dante-cli"]
CMD ["--help"]
