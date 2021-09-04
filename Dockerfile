FROM docker.io/golang:1.17.0@sha256:7dbfeb9d51c049e8bfe36cf1a4217c7b1ba304bf0eb72d57d0c04f405589f122 AS build
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /app/main .
RUN go get github.com/google/go-licenses && go-licenses save ./... --save_path=/notices

FROM gcr.io/distroless/static:nonroot@sha256:be5d77c62dbe7fedfb0a4e5ec2f91078080800ab1f18358e5f31fcc8faa023c4
LABEL org.opencontainers.image.source="https://github.com/greboid/docker-tags-action"
COPY --from=build /notices /notices
COPY --from=build /app/main /docker-tags
WORKDIR /
ENTRYPOINT ["/docker-tags"]
