FROM docker.io/golang:1.17.0@sha256:7dbfeb9d51c049e8bfe36cf1a4217c7b1ba304bf0eb72d57d0c04f405589f122 AS build
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /app/main .
RUN go get github.com/google/go-licenses && go-licenses save ./... --save_path=/notices

FROM gcr.io/distroless/static:nonroot@sha256:c9f9b040044cc23e1088772814532d90adadfa1b86dcba17d07cb567db18dc4e
LABEL org.opencontainers.image.source="https://github.com/greboid/docker-tags-action"
COPY --from=build /notices /notices
COPY --from=build /app/main /docker-tags
WORKDIR /
ENTRYPOINT ["/docker-tags"]
