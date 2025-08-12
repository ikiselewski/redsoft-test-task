FROM golang:latest
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -C cmd/migrate -o migrate
CMD ["./cmd/migrate/migrate"]
RUN go build -C cmd/app -o main .
CMD ["./cmd/app/main"]