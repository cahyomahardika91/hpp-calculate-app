FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /build
COPY . .
RUN go mod tidy
RUN if [ ! -f .env ]; then cp .env.example .env; fi
RUN go build --ldflags "-s -w -extldflags -static" -o main .

FROM alpine:latest

WORKDIR /www

COPY --from=builder /build/main /www/
COPY --from=builder /build/.env /www/.env
COPY --from=builder /build/public/ /www/public/
COPY --from=builder /build/resources/ /www/resources/

EXPOSE 3000

ENTRYPOINT ["/www/main"]
