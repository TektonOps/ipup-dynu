FROM alpine:3.20.2
WORKDIR /ipup
COPY ipup-dynu .
ENTRYPOINT ["/ipup/ipup-dynu"]