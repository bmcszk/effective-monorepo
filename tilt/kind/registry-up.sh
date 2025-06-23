#!/bin/sh

# https://github.com/rpardini/docker-registry-proxy#kind-cluster

mkdir -p $HOME/.cache
mkdir -p $HOME/.cache/docker_registry_proxy
mkdir -p $HOME/.cache/docker_registry_proxy/cache
mkdir -p $HOME/.cache/docker_registry_proxy/certs
echo Cache size: $(du -sh $HOME/.cache/docker_registry_proxy/cache)

(docker network create --ipv6=false kind) || true

docker run -d --rm --name docker_registry_proxy -it \
       --net kind --hostname docker-registry-proxy \
       -p 0.0.0.0:3128:3128 \
       -v $HOME/.cache/docker_registry_proxy/cache:/docker_mirror_cache \
       -v $HOME/.cache/docker_registry_proxy/certs:/ca \
       -e REGISTRIES="quay.io gcr.io ghcr.io artifacts.nvtvt.com" \
       -e ENABLE_MANIFEST_CACHE=true \
       ghcr.io/rpardini/docker-registry-proxy:0.6.4-debug
