FROM golang:1.24-alpine as builder

RUN apk add --no-cache upx tzdata

WORKDIR /src

COPY . .

RUN go mod download

WORKDIR /src/cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app .
RUN upx app

FROM scratch

ARG VERSION
ENV APP_VERSION=$VERSION

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /bin/app

COPY --from=builder /src/cmd/app/app .
COPY config.yaml .

CMD ["./app"]