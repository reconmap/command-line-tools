SHELL := bash
SUBDIRS := $(wildcard */.)

include Golang.mk

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


