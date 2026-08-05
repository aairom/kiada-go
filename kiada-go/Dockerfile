# ---------- build stage ----------
FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY go.mod ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o kiada-go .

# ---------- runtime stage ----------
FROM alpine:3.19
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /build/kiada-go /app/kiada-go
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/kiada-go"]
