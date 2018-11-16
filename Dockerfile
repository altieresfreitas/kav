
FROM golang:1.11 as builder
WORKDIR /go/src/github.com/4ltieres/karepol
COPY . ./
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure -v
RUN CGO_ENABLED=0 GOOS=linux go build -o karepol

FROM scratch

COPY --from=builder /go/src/github.com/4ltieres/karepol/karepol .
CMD ["/karepol"]
