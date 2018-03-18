FROM alpine
MAINTAINER Max Sum <max@lolyculture.com>

# Build app
ADD memreflect
CMD ["./memreflect", "-p", "${MEMREFLECT_PORT:-11211}", "${MEMREFLECT_SHUTDOWN:+-s}"]
