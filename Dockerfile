FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .

RUN go build -ldflags="-s -w" -o /app/billing-engine .


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Jakarta /usr/share/zoneinfo/Asia/Jakarta
ENV TZ Asia/Jakarta

WORKDIR /app
COPY --from=builder /app/billing-engine /app/billing-engine
COPY ./etc /app/etc

CMD ["./billing-engine", "-f", "etc/billing-api.yaml"]
