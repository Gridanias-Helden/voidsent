dirs=$(shell go list -f {{.Dir}} ./...)
mods=$(shell cat go.mod | grep "^module " | head -n 1 - | cut -d" " -f2)

check: test format
	@-

test:
	echo "... Tests"
	@go test -v ./...

format: 
	@echo "... Formats"
	@for d in ${dirs}; do \
		goimports -local ${mods} -l $${d}/*.go; \
		goimports -local ${mods} -w $${d}/*.go; \
	done
