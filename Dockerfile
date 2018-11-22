# Source: https://github.com/rebuy-de/golang-template
# Version: 1.3.1

FROM golang:1.11.2

# Configure Go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

# Install Go Tools
RUN go get -u golang.org/x/lint/golint

# Install Glide
RUN go get -u github.com/Masterminds/glide/...

WORKDIR /go/src/github.com/Masterminds/glide

RUN git checkout v0.12.3
RUN go install

WORKDIR /go/src/github.com/rebuy-de/aws-nuke

CMD make build && chown $UID:$GID aws-nuke && chown -h $UID:$GID aws-nuke
