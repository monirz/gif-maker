FROM golang

WORKDIR /app

ADD . /app

RUN go build main.go

EXPOSE 8090

CMD ["./main"]
