FROM docker.io/library/golang:1.17 AS go-builder
ENV GO111MODULE=off
ENV CGO_ENABLED=1
ENV COMMON_GO_ARGS=-race
ENV GOOS=linux
ENV GOPATH=/go

WORKDIR /go/src/github.com/nwaizer/GetPodLogsEfficiantly
COPY ./main.go ./
RUN go install
CMD ["/go/bin/GetPodLogsEfficiantly"]