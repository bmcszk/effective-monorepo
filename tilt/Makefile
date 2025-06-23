TILT_DIR:=$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

.PHONY: helm-install
helm-install:
	(helm version) || (mkdir -p tmp && cd tmp\
		&& (curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash)\
		&& cd .. && rm -rf tmp\
		&& helm version)

.PHONY: tilt-install
tilt-install:
	(tilt version) || (mkdir -p tmp && cd tmp\
		&& (curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash)\
		&& cd .. && rm -rf tmp\
		&& tilt version)

.PHONY: tilt-up
tilt-up: tilt-install helm-install
	tilt up -f ${TILT_DIR}/Tiltfile

.PHONY: tilt-ci
tilt-ci: tilt-install helm-install
	tilt ci -f ${TILT_DIR}/Tiltfile --timeout 10m

.PHONY: tilt-down
tilt-down:
	tilt down -f ${TILT_DIR}/Tiltfile --delete-namespaces

.PHONY: install-all-deps
install-all-deps: kind-install helm-install tilt-install
	docker version

include ${TILT_DIR}/kind/Makefile
