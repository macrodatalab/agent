FROM golang:1.4.2
MAINTAINER YI-HUNG JEN <yihungjen@macrodatalab.com>

COPY . /go/src/github.com/yihungjen/agent
WORKDIR /go/src/github.com/yihungjen/agent

ENV GOPATH /go/src/github.com/yihungjen/agent/Godeps/_workspace:$GOPATH
RUN go install

ENTRYPOINT ["agent"]
CMD ["--help"]
