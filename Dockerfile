FROM golang:1.12.3 as builder

WORKDIR /crawler
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags '-static'" -o crawler1337 .

FROM debian:9.8-slim as runner

RUN apt-get update && \
    apt-get install --no-install-recommends -y ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /crawler
COPY --from=builder /crawler/crawler1337 .

ENTRYPOINT ["./crawler1337"]
