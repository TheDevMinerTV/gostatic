FROM golang:1.21.0 AS builder
WORKDIR /src

COPY ./go.sum ./go.mod /src/
RUN go mod download

COPY . /src
RUN CGO_ENABLED=0 go build -o /bin/gostatic -ldflags="-w -s"



FROM alpine:3.18.3 AS runner
EXPOSE 3000

COPY --from=builder /bin/gostatic /bin/gostatic

CMD ["/bin/gostatic", "--files", "/static"]
