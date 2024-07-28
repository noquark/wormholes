FROM golang:latest AS compiler
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/wormholes
COPY . .
RUN go mod download

RUN go build -a -installsuffix cgo .

FROM alpine:latest
RUN apk --no-cache add ca-certificates dumb-init
COPY --from=compiler /go/src/wormholes/wormholes /wormholes
EXPOSE 5002
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD [ "/wormholes" ]
