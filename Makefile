REPO_PATH=github.com/echoturing/crypt







.PHONY: build
build:
	@CGO_ENABLED=0 GOOS=darwin go build -o .build/crypt_mac cmd/main.go
	@CGO_ENABLED=0 GOOS=linux go build -o .build/crypt_linux  cmd/main.go
	@CGO_ENABLED=0 GOOS=windows go build -o .build/crypt.exe  cmd/main.go

.PHONY: run
run:
	@go run cmd/main.go

.PHONY: fmt
fmt:
	@find . -name "*.go" | xargs goimports -w -l --local $(REPO_PATH) --private "mockprivate"

.PHONY: vet
vet:
	@go vet ./...

.PHONY: doc
doc:
	@npm run apidoc