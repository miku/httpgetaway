SHELL := /bin/bash
TARGETS := proxytypes.png

.PHONY: all
all: $(TARGETS)

%.png: %.dot
	dot -Tpng $^ > $@

.PHONY: clean
clean:
	rm -f $(TARGETS)
