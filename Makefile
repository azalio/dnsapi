COMMIT=$(shell git rev-parse --short HEAD)
# TAG ?= $(shell git describe --tags --exact-match)


default: ;

make_manifests:
	helm template dnsapi dnsapi-helm/
