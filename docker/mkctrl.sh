#!/bin/bash

docker rmi -f dt/ctrl

ADDGODEPS=`cat ../bootstrap.sh | grep "go *get " | sed -e "s/^/RUN /g"`
if [ $? -ne 0 ]; then
	echo "generate ADDGODEPS failed"
	exit 1
fi

cat > ../Dockerfile << EOF
FROM golang:1.4.2
MAINTAINER zimuxia <zimu_xia@126.com>

RUN apt-get upgrade -y

ENV GOPATH /tmp/go
RUN mkdir -p \${GOPATH}
${ADDGODEPS}

ENV DTPATH /dt
RUN mkdir -p \${DTPATH}

ENV BUILDPATH /tmp/dt
RUN mkdir -p \${BUILDPATH}

ADD pkg \${GOPATH}/src/github.com/pingcap/dt/pkg
ADD cmd \${BUILDPATH}/cmd
WORKDIR \${BUILDPATH}
RUN go build -o \${DTPATH}/ctrl ./cmd
RUN rm -rf \${BUILDPATH}
ADD etc/ctrl_cockroach_cfg.toml \${DTPATH}/

EXPOSE 54320

EOF

docker build -t dt/ctrl ../ && rm -f ../Dockerfile
