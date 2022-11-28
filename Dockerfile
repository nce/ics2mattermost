# syntax=docker.io/docker/dockerfile:1
FROM golang:1.19.1 as build

WORKDIR /go/src
COPY . .
ARG TARGETARCH GO_ARCH
ARG GITSHA VERSION
RUN make build

FROM gcr.io/distroless/static-debian11
COPY --from=build /go/src/ics2mattermost .

USER nobody
CMD ["/ics2mattermost"]
