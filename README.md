# Docker tags action

## Desciption
This action generates a list of docker tags for a given refspec.

- Tags must be semver formatted, and master/main are tagged as latest.  Other refspecs lead to no tags.
- You can specify the image name
- You can specify a list of registries to push to if you're not using docker hub (or want multiples)

## Inputs

 - `repository` - Image name, if you don't want to use the git repository name
 - `registries` - A comma separated list of registry names to push to, defaults to docker hub

## Outputs
 - tag: A comma separated list of tags to be passed to `docker/build-push-action@v2`
 - version: The semver parsed version of the image, or the SHA if this is not tag

## Example

```yaml
name: docker build
on:
  push:
    branches:
      - master
    tags:
      - v*
jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v2
      - name: Generate tags
        id: tags
        uses: greboid/docker-tags-action@v1
      - name: Set up Docker Build
        uses: docker/setup-buildx-action@v1
      - name: Login to DockerHub
        uses: docker/login-action@v1
        with:
          username: ${{ secrets.DOCKER_HUB_USER }}
          password: ${{ secrets.DOCKER_HUB_TOKEN }}
      - name: Build and push
        uses: docker/build-push-action@v2
        with:
          push: true
          tags: ${{ steps.tags.outputs.tags }}
```