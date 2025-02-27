FROM golang:1.23.4-alpine AS builder

RUN apk update && apk --no-cache add bash git make

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /usr/src

COPY ["go.mod","go.sum","./"]

RUN go mod download

COPY . .

# build
RUN go build -o ./bin/app ./cmd/wallet/main.go

FROM alpine:latest AS runner

RUN apk update && apk --no-cache add bash


# copy binary from builder
COPY --from=builder /usr/src/bin/app /app
COPY --from=builder /usr/src/migrations /migrations
COPY --from=builder /usr/src/run-migrations.sh /run-migrations.sh
COPY --from=builder /go/bin/migrate /migrate

RUN chmod +x /run-migrations.sh
RUN chmod +x /migrate

EXPOSE 8080
EXPOSE 8081
EXPOSE 8082

ENTRYPOINT ["/run-migrations.sh"]
CMD ["/app"]