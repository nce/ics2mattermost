# syntax=docker.io/docker/dockerfile:1
FROM golang:1.19.1 as build

WORKDIR /go/src
COPY . .
ARG TARGETARCH TARGETOS
ARG GITSHA VERSION
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    make build

FROM gcr.io/distroless/static-debian11
WORKDIR /
COPY --from=build /go/src/ics2mattermost .

USER nobody
CMD ["/ics2mattermost"]
