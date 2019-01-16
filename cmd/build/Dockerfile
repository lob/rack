FROM golang:1.11

RUN curl -s https://download.docker.com/linux/static/stable/x86_64/docker-18.03.1-ce.tgz | \
    tar -C /usr/bin --strip-components 1 -xz

COPY . $GOPATH/src/github.com/lob/rack

RUN go install github.com/lob/rack/cmd/build
RUN env CGO_ENABLED=0 go install --ldflags '-extldflags "-static"' github.com/lob/rack/cmd/convox-env
