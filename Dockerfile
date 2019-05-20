FROM golang:1.11.10-alpine3.8

ADD https://github.com/golang/dep/releases/download/v0.5.3/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

RUN apk update &&\
    apk add git


RUN mkdir -p /go/src/github.com/billyboar/battleships

WORKDIR /go/src/github.com/billyboar/battleships

COPY . .

RUN dep ensure --vendor-only

EXPOSE 3000

RUN go install

CMD ["battleships"]
