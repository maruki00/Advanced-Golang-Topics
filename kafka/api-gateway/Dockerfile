FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

# RUN go mod tidy
# RUN "go build -o ./api-geteway ./api-geteway"

FROM golang:1.20
COPY --from=builder /app/microservice1 /microservice1
COPY wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh
CMD ["/app/wait-for-it.sh", "kafka:9092", "--", "/microservice1"]
