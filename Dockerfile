FROM golang:1.17.8 as builder
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /go/src/app
COPY . .
RUN go install github.com/cosmtrek/air@v1.29.0
CMD ["air", "-c", ".air.toml"]
