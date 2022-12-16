FROM golang:1.19.3-alpine3.15 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -a -ldflags "-w -s" -o ipinfo-api ./cmd/main.go

FROM alpine
WORKDIR /app 
COPY --from=builder /app/ipinfo-api .
COPY .env .
EXPOSE "${API_PORT}" "${API_PORT}"
CMD [ "./ipinfo-api" ]