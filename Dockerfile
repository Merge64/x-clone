FROM golang:1.23 AS builder

WORKDIR /app

COPY server/go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o x-clone ./server/main/main.go

FROM scratch

WORKDIR /root/

COPY --from=builder /app/x-clone .

CMD ["./x-clone"]