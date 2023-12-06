FROM golang:alpine3.17 as builder
WORKDIR /app
COPY . /app
RUN go build .

FROM golang:alpine3.17 as runtime
COPY --from=builder /app/remote-exec /app/remote-exec
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["./remote-exec"]
