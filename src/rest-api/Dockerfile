FROM golang

WORKDIR /app

COPY . /app

RUN go build cmd/main.go

EXPOSE 8080

CMD ["/app/main"]
