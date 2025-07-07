# build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o cfire .

# run stage
FROM scratch

WORKDIR /app

COPY --from=builder /app/cfire .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/cfire"]


#docker run --rm -e CFIRE_PORT=8888 -p 5000:8888 cfire-monitor
