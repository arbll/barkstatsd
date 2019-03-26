FROM golang:1.10  AS build-env
WORKDIR /
ADD . /go/src/github.com/arbll/barkstatsd
ADD ./build /build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 /build

FROM alpine
COPY --from=build-env /barkstatsd /bark/barkstatsd
ENV BARK_HOST=127.0.0.1 \
    BARK_PORT=8125 \
    BARK_PPS=5000 \
    BARK_INTERVAL=1m \
    BARK_STEP=0 \
    BARK_DURATION=0

VOLUME [ "/input" ]

ENTRYPOINT /bark/barkstatsd -H=$BARK_HOST -p=$BARK_PORT --pps=$BARK_PPS -s=$BARK_STEP -i=$BARK_INTERVAL -d $BARK_DURATION /input
