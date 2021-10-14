FROM ghcr.io/greboid/dockerfiles/golang@sha256:65e504b0cb4e5df85e2301f47cd3f231768d7b0d5aba59b1201e9c50fdf5e0ac AS BUILD

# Build the app
WORKDIR /app
COPY . /app
#Compile the app. Retrieves licenses, set timestamps on the outputs
RUN set -eux; \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -gcflags=./dontoptimizeme=-N -ldflags=-s -o /go/bin/app .; \
    go run github.com/google/go-licenses@latest save ./... --save_path=/notices; \
    mkdir /data; \
    touch --date=@0 /go/bin/app /notices /data

FROM ghcr.io/greboid/dockerfiles/base@sha256:93cb0a11840ca0bf96975b0e9303a234bae7988481533300702d12a5857b630a
LABEL org.opencontainers.image.source="https://github.com/greboid/docker-tags-action"
COPY --from=build /notices /notices
COPY --from=build /go/bin/app /docker-tags
WORKDIR /
ENTRYPOINT ["/docker-tags"]
