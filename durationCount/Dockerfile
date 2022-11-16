FROM golang:1.19 as builder
# ENV DBADDR="db"
WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN mkdir /go/bin/app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /go/bin/app -v ./...

FROM alpine:latest as run
RUN apk --no-cache add ca-certificates
RUN addgroup -S app && adduser -S app -G app
COPY --from=builder --chown=app /go/bin/app /app
RUN chmod +x /app/*
USER app

CMD [ "/app/durationCount" ]