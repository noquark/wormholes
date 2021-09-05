FROM golang:latest AS compiler
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/wormholes
COPY go.* .
RUN go mod download
COPY . .
RUN go build -a -installsuffix cgo ./services/director
RUN go build -a -installsuffix cgo ./services/generator
RUN go build -a -installsuffix cgo ./services/creator

FROM alpine:latest as director
RUN apk --no-cache add ca-certificates
COPY --from=compiler /go/src/wormholes/director /
EXPOSE 5000
CMD /director

FROM alpine:latest as generator
RUN apk --no-cache add ca-certificates
COPY --from=compiler /go/src/wormholes/generator /
EXPOSE 5001
CMD /generator

FROM alpine:latest as creator
RUN apk --no-cache add ca-certificates
COPY --from=compiler /go/src/wormholes/creator /
EXPOSE 5002
CMD /creator
