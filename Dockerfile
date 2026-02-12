FROM --platform=linux/amd64 golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY back ./back
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/todo ./cmd

FROM alpine:latest

WORKDIR /app
COPY --from=builder /out/todo /app/todo
COPY web /app/web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db

EXPOSE 7540

CMD ["/app/todo"]