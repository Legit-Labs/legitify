FROM golang:1.18.2-alpine as build_env

WORKDIR /

ARG OS
ARG ARCH

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build \
    -ldflags "-s -w" \
    -a -installsuffix cgo \
    -o legitify

FROM scratch

COPY --from=build_env /legitify /legitify

WORKDIR /
ENTRYPOINT ["/legitify"]