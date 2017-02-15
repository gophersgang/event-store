FROM golang:1.7
RUN mkdir -p /go/src/app
WORKDIR /go/src/app

RUN go get -u github.com/golang/lint/golint

COPY . /go/src/github.com/vendasta/event-store

RUN go build -o /bin/event-store github.com/vendasta/event-store/server/
CMD ["/bin/event-store"]
