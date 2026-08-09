FROM golang:1.26-alpine
RUN apk add git

WORKDIR /plum

COPY ./plum-backend ./
RUN go install -v ./cmd/plum

CMD ["/bin/sh", "-c", "plum run --verbosity debug /config.json"]
