FROM golang:1.22.3-bullseye as builder 

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify 

COPY . .

RUN CGO_ENABLED=0 go build -o ./main ./cmd/main.go

FROM alpine:3.19 AS executer 

COPY --from=builder /app/main / 

EXPOSE 8020 

CMD ["/main"]
