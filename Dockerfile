FROM golang:1.26-alpine AS builder

WORKDIR /app

# we do this first to separate cache layers for the modules and the code
# if code changes, that would invalide the cache and we'd have to download all modules again
# this way, we can have a cache just for modules
COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /out/controller ./cmd/controller/

FROM alpine:3.24

COPY --from=builder /out/controller /bin/controller

CMD ["/bin/controller"]

EXPOSE 9000


# TODO: Move to compose file

# saving for broker env var
# ENV BROKER_ID="broker-1"
# ENV BROKER_ADDR="localhost:9002"
# ENV BROKER_PORT="9002"
# ENV BROKER_PORT="9002"
