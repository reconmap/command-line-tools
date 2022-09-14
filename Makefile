SHELL := bash
SUBDIRS := $(wildcard */.)

all: $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all $(SUBDIRS)

clean:
	pushd agent && make clean && popd
	pushd cli && make clean && popd

