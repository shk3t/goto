FROM golang:1.22-alpine
WORKDIR /go-redis
COPY go.mod go.sum .
RUN go mod download
COPY src src
CMD ["go", "test", "-count=1", "-v", "./..."]
