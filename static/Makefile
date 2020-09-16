SHELL := /bin/bash
TARGETS := proxytypes.png fluxproxy.png

.PHONY: all
all: $(TARGETS)

%.png: %.dot
	dot -Tpng $^ > $@

.PHONY: clean
clean:
	rm -f $(TARGETS)
