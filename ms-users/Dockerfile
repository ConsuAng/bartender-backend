FROM golang:1.22.2 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]