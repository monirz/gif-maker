FROM golang:1.12

ENV GO111MODULE=on
ENV PORT=8090
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

EXPOSE 8090
ENTRYPOINT ["/app/gif-maker"]