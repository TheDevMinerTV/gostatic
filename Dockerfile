FROM golang:1.23.3 AS builder
WORKDIR /src

COPY ./go.sum ./go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /bin/gostatic -ldflags="-w -s"

FROM alpine:3.20 AS runner
RUN adduser -D -u 1000 app && \
    apk add --no-cache mailcap
COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE 80

COPY --from=builder /bin/gostatic /bin/gostatic

ENTRYPOINT ["/entrypoint.sh"]
