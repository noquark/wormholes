FROM golang:latest AS compiler
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/wormholes
COPY . .
RUN go mod download

RUN go build -a -installsuffix cgo ./services/creator
RUN go build -a -installsuffix cgo ./services/generator
RUN go build -a -installsuffix cgo ./services/redirector

FROM alpine:latest as creator
RUN apk --no-cache add ca-certificates dumb-init
COPY --from=compiler /go/src/wormholes/creator /creator
EXPOSE 5002
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD [ "/creator" ]

FROM alpine:latest as generator
RUN apk --no-cache add ca-certificates dumb-init
COPY --from=compiler /go/src/wormholes/generator /generator
EXPOSE 5001
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD [ "/generator" ]

FROM alpine:latest as redirector
RUN apk --no-cache add ca-certificates dumb-init
COPY --from=compiler /go/src/wormholes/redirector /redirector
EXPOSE 5000
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD [ "/redirector" ]
