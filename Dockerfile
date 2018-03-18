FROM golang:alpine
MAINTAINER Max Sum <max@lolyculture.com>

# Build app

RUN apk add --no-cache git iptables \
    && go get -t github.com/Max-Sum/memreflect/build \
    && apk del git \
    && go build -o memreflect github.com/Max-Sum/memreflect/build

CMD ["./memreflect", "-p", "$MEMREFLECT_PORT", "${MEMREFLECT_SHUTDOWN:+-s}"]
