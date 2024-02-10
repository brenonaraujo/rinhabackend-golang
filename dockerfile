FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest  

RUN apk --no-cache add ca-certificates netcat-openbsd

WORKDIR /root/

COPY --from=builder /app/main .

COPY db_wait.sh /root/db_wait.sh

RUN chmod +x /root/db_wait.sh

EXPOSE 33888

ENTRYPOINT ["/root/db_wait.sh", "db", "5432"]

CMD ["./main"]
