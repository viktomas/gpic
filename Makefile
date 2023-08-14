SHELL = /bin/bash
.SILENT: # comment this one out if you need to debug the Makefile
APP=picsort

.PHONY: run
run:
	go run . ~/Downloads/tmp/pictures/

.PHONY: build
build:
	go build
	# go build -o ${APP} main.go

.PHONY: test
test:
	go test ./...

.PHONY: watch-test
watch-test:
	fswatch \
		--exclude '.git' \
		-o . | \
		xargs -n1 -I{} go test ./...

.PHONY: clean
clean:
	go clean


.PHONY: watch-run
watch-run:
	fswatch \
		--exclude '.git' \
		-xn . | \
		while read file event; do \
			echo "File $$file has changed, Event: $$event"; \
			$(MAKE) run ; \
		done
