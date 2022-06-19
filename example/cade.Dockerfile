FROM alpine:latest

RUN apk add --no-cache bash

RUN adduser -Ss /bin/bash cadeuser

RUN chown -hR cadeuser: /home/cadeuser

WORKDIR /home/cadeuser

USER cadeuser

RUN mkdir -p workdir

WORKDIR /home/cadeuser/workdir

RUN echo "hello from a container workspace!" > hello.txt