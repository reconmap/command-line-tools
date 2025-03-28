SHELL := bash
SUBDIRS := $(wildcard */.)

all: $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all $(SUBDIRS)

programs:
	pushd agent && make reconmapd && popd
	pushd cli && make rmap && popd

clean:
	pushd agent && make clean && popd
	pushd cli && make clean && popd

.PHONY: lint
lint: GOLANGCI_LINT_VERSION ?= 2.0.2
lint:
	docker run \
	-v $(CURDIR):/reconmap/agent \
	-w /reconmap/agent \
	golangci/golangci-lint:v$(GOLANGCI_LINT_VERSION)-alpine \
	golangci-lint run -c .golangci.yml --timeout 10m --fix agent


