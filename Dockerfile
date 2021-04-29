FROM golang:latest AS compiler
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/wormholes
COPY go.* .
RUN go mod download
COPY . .
RUN go build -a -installsuffix cgo .

FROM alpine:latest as runner
COPY --from=compiler /go/src/wormholes/wormholes /
EXPOSE 3000
CMD ["/wormholes"]
