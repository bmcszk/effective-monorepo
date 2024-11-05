KIND_DIR:=$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

.PHONY: kind-install
kind-install:
	(kind version) || (go install sigs.k8s.io/kind@latest\
		&& kind version)

.PHONY: docker-registry-proxy-up
docker-registry-proxy-up:
	${KIND_DIR}/registry-up.sh

.PHONY: docker-registry-proxy-down
docker-registry-proxy-down:
	${KIND_DIR}/registry-down.sh

.PHONY: kind-up
kind-up: kind-install docker-registry-proxy-up
	kind create cluster --config ${KIND_DIR}/cluster.yaml

.PHONY: kind-down
kind-down: docker-registry-proxy-down
	kind delete cluster
