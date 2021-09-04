FROM docker.io/golang:1.17@sha256:33ef0040801bb4deabe1db381ee92de1afc81b869ce27d52fb52d24cf37ff543 AS build
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
