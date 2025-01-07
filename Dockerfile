FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o rpl-service ./main/main.go

FROM scratch

WORKDIR /root/

COPY --from=builder /app/rpl-service .

CMD ["./rpl-service"]