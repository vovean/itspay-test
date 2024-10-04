FROM golang:1.23-alpine3.20 AS builder

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app

#############################
FROM alpine:3.20 AS runner

# Use non-root user for security
RUN adduser -D -h /home/runner runner
USER runner
WORKDIR /home/runner

COPY --from=builder --chown=runner:runner /bin/app ./app
COPY config ./config