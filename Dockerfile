FROM alpine:3.16.2 as certs
RUN apk --update add ca-certificates=20220614-r0

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY legitify /legitify
ENTRYPOINT ["/legitify"]