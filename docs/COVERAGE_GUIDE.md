1. Verifique a seção [Troubleshooting](#troubleshooting)
2. Consulte o [README](./pkg/coverage/README.md)
3. Abra uma issue no repositório
# Guia de Integração - Coverage Report Generator

Documentação detalhada sobre como usar a biblioteca de geração de relatórios de cobertura.

## 📋 Sumário

1. [Instalação e Setup](#instalação-e-setup)
2. [Uso Básico](#uso-básico)
3. [Uso como Biblioteca](#uso-como-biblioteca)
4. [Uso como CLI](#uso-como-cli)
5. [Integração com CI/CD](#integração-com-cicd)
6. [Integração com Makefile](#integração-com-makefile)
7. [Exemplos Avançados](#exemplos-avançados)
8. [Troubleshooting](#troubleshooting)

## Instalação e Setup

### Pré-requisitos
- Go 1.20+
- Git

### Estrutura de Diretórios
```
projeto/
├── pkg/
│   └── coverage/          # Biblioteca
│       ├── parser.go
│       ├── html_generator.go
│       ├── parser_test.go
│       └── README.md
├── cmd/
│   └── coverage-report/   # CLI
│       └── main.go
├── scripts/
│   └── generate-coverage-report.sh  # Script auxiliar
└── Makefile               # Targets de build
```

## Uso Básico

### Como Biblioteca

```go
package main

import (
	"log"
	"os"
	"shipping-management/pkg/coverage"
)

func main() {
	// 1. Abrir arquivo de cobertura
	file, err := os.Open("coverage.out")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 2. Parsear cobertura
	cov, err := coverage.ParseCoverageFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Acessar informações
	fmt.Printf("Total de cobertura: %.2f%%\n", cov.GetTotalCoverage())
	fmt.Printf("Modo: %s\n", cov.Mode)
	fmt.Printf("Arquivos: %d\n", len(cov.Files))

	// 4. Gerar relatório HTML
	generator := coverage.NewHTMLGenerator(cov)
	
	output, err := os.Create("report.html")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	err = generator.Generate(output)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Uso como CLI

### Opções Básicas

```bash
# Usando valores padrão
go run cmd/coverage-report/main.go

# Especificar entrada
go run cmd/coverage-report/main.go -in custom-coverage.out

# Especificar entrada e saída
go run cmd/coverage-report/main.go -in coverage.out -out my-report.html

# Versão compilada
./coverage-report -in coverage.out -out report.html
```

### Compilação

```bash
# Compilar CLI
go build -o coverage-report cmd/coverage-report/main.go

# Instalar globalmente
go install ./cmd/coverage-report
```

## Integração com CI/CD

### GitHub Actions

Criar arquivo `.github/workflows/coverage.yml`:

```yaml
name: Coverage Report

on: [push, pull_request]

jobs:
  coverage:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Run tests with coverage
        run: |
          go test -coverpkg=./... -coverprofile=coverage.out ./...
      
      - name: Generate coverage report
        run: |
          go run cmd/coverage-report/main.go -in coverage.out -out coverage-report.html
      
      - name: Upload coverage report
        uses: actions/upload-artifact@v3
        with:
          name: coverage-report
          path: coverage-report.html
```

### GitLab CI

Criar arquivo `.gitlab-ci.yml`:

```yaml
coverage:
  image: golang:1.24
  stage: test
  script:
    - go test -coverpkg=./... -coverprofile=coverage.out ./...
    - go run cmd/coverage-report/main.go -in coverage.out -out coverage-report.html
  artifacts:
    paths:
      - coverage-report.html
    expire_in: 30 days
```

## Integração com Makefile

Adicionar targets ao seu Makefile:

```makefile
# test: Executa testes e gera cobertura padrão
test:
	go test -race -coverpkg=./internal/... -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# coverage-report: Gera relatório HTML interativo
coverage-report:
	@echo "🎨 Gerando relatório de cobertura interativo..."
	go run cmd/coverage-report/main.go -in coverage.out -out coverage-report.html
	@echo "✨ Relatório gerado: coverage-report.html"

# test-report: Executa testes e gera relatório HTML
test-report: test coverage-report
	@echo "📊 Relatório disponível em coverage-report.html"

# open-report: Abre o relatório no navegador
open-report: coverage-report
	@command -v xdg-open >/dev/null 2>&1 && xdg-open coverage-report.html || \
	 command -v open >/dev/null 2>&1 && open coverage-report.html || \
	 echo "Abra coverage-report.html manualmente"

.PHONY: test coverage-report test-report open-report
```

Uso:

```bash
make test-report        # Testa e gera relatório
make open-report        # Abre relatório no navegador
make coverage-report    # Apenas gera relatório
```

## Exemplos Avançados

### Gerar Múltiplos Relatórios

```bash
#!/bin/bash

# Gerar cobertura para diferentes pacotes
for pkg in "./internal/application" "./internal/infrastructure" "./internal/domain"; do
  output="coverage-$(basename $pkg).out"
  go test -coverprofile="$output" "$pkg/..." 
  go run cmd/coverage-report/main.go -in "$output" -out "report-$(basename $pkg).html"
done

echo "Relatórios gerados:"
ls -lh report-*.html
```

### Comparar Cobertura

```bash
#!/bin/bash

# Gerar dois relatórios para comparação
echo "Cobertura main:"
go test -coverprofile=coverage-main.out ./...
go tool cover -func=coverage-main.out | tail -1

echo ""
echo "Cobertura develop:"
git stash
go test -coverprofile=coverage-develop.out ./...
go tool cover -func=coverage-develop.out | tail -1
git stash pop

# Gerar relatórios
go run cmd/coverage-report/main.go -in coverage-main.out -out report-main.html
go run cmd/coverage-report/main.go -in coverage-develop.out -out report-develop.html
```

### Filtro de Cobertura Mínima

```go
package main

import (
	"log"
	"os"
	"shipping-management/pkg/coverage"
)

func main() {
	file, _ := os.Open("coverage.out")
	defer file.Close()

	cov, _ := coverage.ParseCoverageFile(file)
	total := cov.GetTotalCoverage()

	// Falhar se cobertura < 80%
	if total < 80 {
		log.Fatalf("Cobertura insuficiente: %.2f%% (mínimo: 80%%)", total)
	}

	log.Printf("✅ Cobertura aceitável: %.2f%%", total)
}
```

## Troubleshooting

### Problema: "arquivo de cobertura não encontrado"

**Solução:**
```bash
# Gerar primeiro o arquivo de cobertura
go test -coverpkg=./... -coverprofile=coverage.out ./...

# Depois gerar o relatório
go run cmd/coverage-report/main.go
```

### Problema: "relatório não mostra todos os arquivos"

**Solução:**
Especifique o pacote correto ao rodar testes:
```bash
go test -coverpkg=./... -coverprofile=coverage.out ./...
#       ^^^^^^^^^^
#       Importante: incluir todos os pacotes
```

### Problema: "HTML não funciona no navegador"

**Solução:**
- Verificar se o arquivo não está truncado: `wc -l coverage-report.html`
- Abrir em outro navegador
- Verificar console de erro (F12)

### Problema: "Cobertura mostra 0% para alguns arquivos"

**Causa:** Arquivo não está sendo testado

**Solução:**
```bash
# Verificar quais arquivos estão sendo testados
go test -coverpkg=./... -v ./...

# Ajustar -coverpkg para incluir os pacotes desejados
```

## Performance

### Tempos Aproximados

| Operação | Tempo |
|----------|-------|
| Parse de 5 arquivos | < 10ms |
| Parse de 100 arquivos | < 50ms |
| Geração de HTML | 50-200ms |
| Total (5 arquivos) | < 300ms |

### Tamanho do Relatório

| Quantidade | Tamanho |
|-----------|---------|
| 5 arquivos | ~12KB |
| 50 arquivos | ~50KB |
| 100+ arquivos | ~150KB |

## Referência Rápida

```bash
# Gerar cobertura
go test -coverpkg=./... -coverprofile=coverage.out ./...

# Ver resumo no terminal
go tool cover -func=coverage.out

# Gerar relatório interativo
go run cmd/coverage-report/main.go

# Compilar CLI
go build -o coverage-report cmd/coverage-report/main.go

# Usar Makefile
make test-report
make open-report
```

## Recursos Adicionais

- [Documentação do Go Coverage](https://go.dev/blog/cover)
- [README da Biblioteca](./pkg/coverage/README.md)
- [Código Fonte](./pkg/coverage/)

## Suporte

Para dúvidas ou problemas:

