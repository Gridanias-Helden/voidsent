dirs=$(shell go list -f {{.Dir}} ./...)
mods=$(shell cat go.mod | grep "^module " | head -n 1 - | cut -d" " -f2)

check: test format
	@-

test:
	echo "... Tests"
	@go test -v ./...

format: 
	@echo "... Formats"
	@echo ${dirs}
	@for d in ${dirs}; do \
		echo $${d}; \
		goimports -local ${mods} -l $${d}/*.go; \
		goimports -local ${mods} -w $${d}/*.go; \
	done
