FROM "golang:alpine" as builder
WORKDIR $GOPATH/src/build
COPY . $GOPATH/src/build/
RUN apk add git && go get -v ./... &&  go build -o main

FROM alpine
WORKDIR /app
COPY --from=builder /go/src/build/main /app/main
EXPOSE 8000
ENTRYPOINT ["./main"]