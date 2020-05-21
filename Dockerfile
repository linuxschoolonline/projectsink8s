FROM "golang:alpine" as builder
WORKDIR $GOPATH/src/build
COPY . $GOPATH/src/build/
RUN apk add git && go get -v ./... &&  go build -o main

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /go/src/build/static /app/static
COPY --from=builder /go/src/build/templates /app/templates
COPY --from=builder /go/src/build/main /app/main
EXPOSE 8080
WORKDIR /app
ENTRYPOINT ["./main"]