# Image name: msw-scrapper
FROM golang

ADD . /go/src/ridenow/msw/

# TODO: vendoring
RUN go get github.com/PuerkitoBio/goquery
RUN go get github.com/lib/pq
RUN go get github.com/streadway/amqp

RUN go build -o scrapper /go/src/ridenow/msw/cmd/main.go

# TODO: this should be in /go/bin probably
ENTRYPOINT /go/src/ridenow/msw/scrapper

EXPOSE 8080