FROM golang:1.26.3

WORKDIR /app

COPY shared/ ./shared/
COPY payment-service/go.mod payment-service/go.sum ./
RUN go mod edit -replace shared=./shared
RUN go mod download

COPY payment-service/ .

RUN go mod edit -replace shared=./shared
RUN go build -o payment-service ./main.go

CMD ["./payment-service"]