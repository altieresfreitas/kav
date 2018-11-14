
FROM golang:1.11 as builder
WORKDIR /go/src/github.com/4ltieres/k8s-opol
COPY . ./
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure -v
RUN CGO_ENABLED=0 GOOS=linux go build -o k8s-opol

FROM scratch

COPY --from=builder /go/src/github.com/4ltieres/k8s-opol/k8s-opol .
CMD ["/k8s-opol"]
