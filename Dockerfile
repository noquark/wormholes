FROM golang:latest AS compiler
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/wh

COPY go.* ./
RUN go mod download

COPY . /go/src/wh
RUN go build -a -installsuffix cgo -o wh ./cmd/wh

FROM alpine:latest as runner
WORKDIR /wh
COPY --from=compiler /go/src/wh/wh .
EXPOSE 3000
CMD ["./wh"]
