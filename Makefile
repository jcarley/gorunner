NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

deps:
	@echo "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	@go get -d -v ./...
	@echo $(DEPS) | xargs -n1 go get -d
	@go get github.com/onsi/ginkgo/ginkgo
	@go get github.com/onsi/gomega

updatedeps:
	@echo "$(OK_COLOR)==> Updating all dependencies$(NO_COLOR)"
	@go get -d -v -u ./...
	@echo $(DEPS) | xargs -n1 go get -d -u
	@go get -d -v -u github.com/onsi/ginkgo/ginkgo
	@go get -d -v -u github.com/onsi/gomega

format:
	go fmt ./...

test: deps
	@echo "$(OK_COLOR)==> Testing gorunner...$(NO_COLOR)"
	go test ./...
