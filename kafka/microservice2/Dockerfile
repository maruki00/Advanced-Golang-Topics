FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

# RUN go mod tidy
# RUN "go build -o ./microservice2 ./microservice2"

FROM golang:1.20
COPY --from=builder /app/microservice2 /microservice2
COPY wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh
CMD ["/app/wait-for-it.sh", "kafka:9092", "--", "/microservice2"]
