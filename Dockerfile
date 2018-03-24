FROM golang:alpine AS build-env
LABEL maintainer "Max Sum <max@lolyculture.com>"

# Build app

RUN apk add --update git \
    && go get -t github.com/Max-Sum/memreflect/build \
    && go build -o /memreflect github.com/Max-Sum/memreflect/build

# Final stage
FROM alpine

RUN apk add --no-cache iptables
COPY --from=build-env /memreflect /

CMD /memreflect -p ${MEMREFLECT_PORT:-11211} ${MEMREFLECT_SHUTDOWN:+-s}
