# Makefile para Coverage Report Generator

.PHONY: test build install clean help

# test: Executar testes
test:
	go test -race -v -cover ./...

# build: Compilar CLI
build:
	go build -o coverage-report ./cmd/coverage-report

# install: Instalar CLI globalmente
install:
	go install ./cmd/coverage-report

# test-coverage: Gerar arquivo de cobertura
test-coverage:
	go test -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# example: Gerar exemplo de relatório
example: test-coverage
	go run ./cmd/coverage-report -in coverage.out -out example-report.html
	@echo "✅ Relatório de exemplo gerado: example-report.html"

# clean: Limpar arquivos gerados
clean:
	rm -f coverage.out coverage-report.html example-report.html
	rm -f coverage-report
	go clean

# help: Mostrar esta mensagem
help:
	@grep -E '^# [a-z]' Makefile | sed 's/^# /  make /'

.DEFAULT_GOAL := help

