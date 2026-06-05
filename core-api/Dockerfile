FROM golang:1.26.3

WORKDIR /app

COPY shared/ ./shared/
COPY core-api/go.mod core-api/go.sum ./
RUN go mod edit -replace shared=./shared
RUN go mod download

COPY core-api/ .

RUN go mod edit -replace shared=./shared
RUN go build -o servidor ./cmd

EXPOSE 8080

CMD ["./servidor"]