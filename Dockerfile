FROM golang:1.23-alpine AS builder

RUN apk add --no-cache upx
RUN apk --no-cache add tzdata

WORKDIR /src

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o crypto-trade main.go
RUN upx crypto-trade


FROM scratch

# take env from build args
ARG VERSION
ENV APP_VERSION=$VERSION

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /bin/crypto-trade

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /src/crypto-trade .

# Copy the config file
COPY config.yaml /bin/crypto-trade/

CMD ["./crypto-trade"]
