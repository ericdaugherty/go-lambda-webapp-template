GOOS?=$(shell go env GOOS)
RACE?=-race

.PHONY: precompile
precompile: ## Runs pre-compile steps for deployment
	pkger

.PHONY: build
build: ## Builds go application
	env GOOS=$(GOOS) go build $(RACE) -o webapp

.PHONY: run
run: build ## Runs the app locally
	./webapp

.PHONY: clean
clean: ## Cleanup
	rm -f webapp

.PHONY: aws
aws: 
	$(eval GOOS=linux)
	$(eval RACE=)

.PHONY: log
log: ## Tail Lambda logfile
	serverless logs -f web -t

.PHONY: deploy
deploy: clean aws precompile build ## Deploy via Serverless
	sls deploy --verbose

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.DEFAULT_GOAL := help