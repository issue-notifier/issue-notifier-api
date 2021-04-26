FROM golang:alpine3.13

WORKDIR /app

COPY . .

RUN go build -o main .
CMD ["/app/main"]

EXPOSE 8001