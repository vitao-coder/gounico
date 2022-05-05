FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

ADD . /gounico/config

ADD . /gounico/logs/

RUN touch /gounico/logs/gounico.log

COPY /config/config.docker.yaml /gounico/config/config.yaml

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .

EXPOSE 8008

EXPOSE 8009

CMD ["/dist/main"]