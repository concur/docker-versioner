FROM alpine:3.7

RUN apk add --update git

COPY publish/versioner /usr/bin/versioner

ENTRYPOINT ["/usr/bin/versioner"]