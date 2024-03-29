FROM golang:1.22 AS builder

LABEL authors="log"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

ARG SECRET
ENV SECRET=$SECRET

ARG ENV
ENV ENV=$ENV

ARG TOKEN
ENV TOKEN=$TOKEN

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 80

CMD ["./main"]
