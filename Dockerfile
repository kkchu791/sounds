FROM golang:1.26-alpine AS builder

WORKDIR /app

# we do this first to separate cache layers for the modules and the code
# if code changes, that would invalide the cache and we'd have to download all modules again
# this way, we can have a cache just for modules
COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /out/controller ./cmd/controller/ && \
    CGO_ENABLED=0 go build -o /out/broker ./cmd/broker

FROM alpine:3.24 AS controller
COPY --from=builder /out/controller /bin/controller
EXPOSE 9000
CMD ["/bin/controller"]

FROM alpine:3.24 AS broker
COPY --from=builder /out/broker /bin/broker
CMD ["/bin/broker"]
EXPOSE 9001

