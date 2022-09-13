FROM golang:1.17-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/

ENV CGO_ENABLED=0
RUN go build -o /build/reconmapd ./cmd/reconmapd

FROM scratch

WORKDIR /app

COPY --from=builder /build/reconmapd /app/reconmapd

EXPOSE 5520

CMD ["/app/reconmapd"]

