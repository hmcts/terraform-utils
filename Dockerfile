# Build env
FROM golang:1.12 as build

WORKDIR /go/src/tf-utils
COPY . .

RUN go get -d -v ./...
RUN go test -v ./...
RUN go install -v ./...


# Runtime env
FROM gcr.io/distroless/base:debug
COPY --from=build /go/bin/tf-utils /
ENTRYPOINT ["/tf-utils"]
