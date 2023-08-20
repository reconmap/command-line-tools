FROM golang:1.18-alpine AS builder

WORKDIR /build

RUN mkdir -p agent/{cmd,internal}
COPY agent/go.mod agent/go.sum ./agent
COPY agent/cmd/ ./agent/cmd
COPY agent/internal/ ./agent/internal
COPY shared-lib/ ./shared-lib

WORKDIR /build/agent
ENV CGO_ENABLED=0
RUN go build -o /build/reconmapd ./cmd/reconmapd

FROM kalilinux/kali-rolling

WORKDIR /app

COPY --from=builder /build/reconmapd /app/reconmapd

EXPOSE 5520

CMD ["/app/reconmapd"]

