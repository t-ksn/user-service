FROM alpine:3.5
MAINTAINER Tochka <vinogradovve@gmail.com>
RUN apk add --update ca-certificates
ARG VERSION=unkown
LABEL VERSION=$VERSION
COPY app .

EXPOSE 8080
ENTRYPOINT ["/app"]
