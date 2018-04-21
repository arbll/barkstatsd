FROM golang:1.10.1  AS build-env
WORKDIR /
ADD . /go/src/github.com/arbll/barkstatsd
ADD ./build /build
RUN /build

FROM alpine
WORKDIR /bark
COPY --from=build-env /barkstatsd /bark/
ENTRYPOINT ./barkstatsd