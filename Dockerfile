FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY content ./content
COPY templates ./templates
COPY static ./static

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/realtek-connect ./cmd/server

FROM alpine:3.22

WORKDIR /app

RUN addgroup -S app && adduser -S -G app app \
    && mkdir -p /data \
    && chown -R app:app /app /data

COPY --from=builder /out/realtek-connect /app/realtek-connect
COPY --from=builder /src/content /app/content
COPY --from=builder /src/templates /app/templates
COPY --from=builder /src/static /app/static

ENV PORT=8080
ENV DATABASE_PATH=/data/connectplus.db
ENV ANALYTICS_DATABASE_PATH=/data/analytics.db

VOLUME ["/data"]
EXPOSE 8080

USER app

ENTRYPOINT ["/app/realtek-connect"]
