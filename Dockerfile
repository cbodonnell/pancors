# Build stage
FROM golang:1-alpine as builder

RUN mkdir /pancors
WORKDIR /pancors

COPY go.mod go.mod
COPY pancors.go pancors.go
COPY cmd/ cmd/

RUN go get -d -v ./...
RUN go build -o pancors ./cmd/pancors

# Production stage
FROM alpine:latest as prod

RUN mkdir /pancors
WORKDIR /pancors

COPY --from=builder /pancors/pancors ./

EXPOSE 8080

CMD [ "./pancors" ]