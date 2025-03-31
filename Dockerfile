FROM golang:1.24-alpine3.21 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /server ./cmd/main.go

FROM alpine:3.21 AS run

WORKDIR /app

COPY --from=build /server /server

EXPOSE 3333

CMD ["/server"]