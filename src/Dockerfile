FROM golang:1.11.1-alpine3.8 as base

RUN set -eux; \
	apk add --no-cache --virtual git

RUN apk add --update \
        bash \
        ca-certificates \
        tzdata

WORKDIR /go/src/app
COPY . .

RUN go get .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o main

FROM scratch
COPY --from=base /go/src/app/main /distil-proof-of-nothing
CMD ["/distil-proof-of-nothing"]
EXPOSE 3000