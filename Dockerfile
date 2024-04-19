FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o main ./cmd
 
FROM alpine:latest
RUN apk update
RUN apk add --no-cache bash curl jq
WORKDIR /app

ENV DATABASE_URL=mongodb://host.docker.internal:27017
ENV DATABASE_NAME=nub
ENV JWT_SECRET=H@rleyd@vids0n
ENV JWT_EXP=24

COPY --from=builder /app/main ./
CMD ["./main"]