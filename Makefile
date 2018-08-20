all: build

clean:
	@echo Cleaning...
	@rm -rf build

build: clean
	@./scripts/build

test:
	@echo Running unit tests...
	@./scripts/test

check: test

acceptance-test: build
	@echo Running acceptance tests...
	@./scripts/acceptance-test

lint:
	@echo Running linter...
	@./scripts/lint

ci:
	@./scripts/ci
