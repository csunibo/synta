FROM golang:alpine AS go-builder
COPY . /build
WORKDIR /build/cmd
RUN go build -ldflags "-s -w" -o /build/synta

FROM alpine
COPY --from=go-builder /build/synta /usr/bin/synta

ENTRYPOINT ["/usr/bin/synta"]
