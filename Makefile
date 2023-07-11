dirs=$(shell go list -f {{.Dir}} ./...)
mods=$(shell cat go.mod | grep "^module " | head -n 1 - | cut -d" " -f2)

check: format test
	@-

test:
	echo "... Tests"
	@go test -v ./...

fmt: 
	@echo "... Formats"
	@for d in ${dirs}; do \
		goimports -l -local ${mods} -w $${d}/*.go; \
	done
