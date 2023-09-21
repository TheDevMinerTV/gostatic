FROM golang:1.21.1 AS builder
WORKDIR /src

COPY ./go.sum ./go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /bin/gostatic -ldflags="-w -s"

FROM alpine:3.18.3 AS runner
RUN adduser -D -u 1000 app && \
    apk add --no-cache mailcap
EXPOSE 80

COPY --from=builder /bin/gostatic /bin/gostatic

USER app

ENTRYPOINT ["/bin/gostatic", "--files", "/static", "--addr", ":80"]
