# Start from the latest golang base image
FROM golang:1.18.4-alpine3.16 as builder

COPY . /app/smtp-proxy

# build
RUN apk add --no-cache git build-base linux-headers && \
    cd /app/smtp-proxy && go mod tidy && CGO_ENABLED=0 make

######## Start a new stage #######
FROM alpine:3.16

LABEL maintainer="none<none.one>"

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/smtp-proxy/smtp-proxy /bin/smtp-proxy

RUN apk upgrade --update && chmod +x /bin/smtp-proxy

EXPOSE 25
WORKDIR /

ENTRYPOINT ["/bin/smtp-proxy"]
