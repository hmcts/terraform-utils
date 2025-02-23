# Build env
FROM golang:1.24 as build

WORKDIR /go/src/tf-utils
COPY . .

ENV GO111MODULE on
RUN go get -d -v ./...
RUN go test -v ./...
RUN go install -v ./...


# Runtime env
FROM alpine:3.21
COPY --from=build /go/bin/tf-utils /
CMD ["/tf-utils"]
