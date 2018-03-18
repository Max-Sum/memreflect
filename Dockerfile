FROM alpine
MAINTAINER Max Sum <max@lolyculture.com>

ADD memreflect .
RUN apk add --no-cache iptables
CMD "./memreflect -p ${MEMREFLECT_PORT:-11211} ${MEMREFLECT_SHUTDOWN:+-s}"
