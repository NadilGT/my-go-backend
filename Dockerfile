FROM golang:1.24-alpine

WORKDIR /app

COPY go.* ./

RUN go mod tidy

COPY . .

RUN go build -o main main.go
# RUN go build -o /TicketMgt (build all the files)

EXPOSE 8080

CMD ["./main"]