FROM fedora:latest

RUN dnf install -y go

RUN dnf module install -y nodejs:18/common

RUN useradd -ms /bin/bash cadeuser

RUN chown -hR cadeuser: /home/cadeuser

WORKDIR /home/cadeuser

USER cadeuser

RUN mkdir -p workdir

WORKDIR /home/cadeuser/workdir