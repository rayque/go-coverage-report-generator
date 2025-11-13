# 🚀 Quick Start - Coverage Report Generator

Comece em **5 minutos**!

## Instalação

### Opção 1: Usar como CLI (Mais fácil)

```bash
# Instalar globalmente
go install github.com/rayque.oliveira/coverage-report-generator/cmd/coverage-report@latest

# Usar
go test -coverpkg=./... -coverprofile=coverage.out ./...
coverage-report -in coverage.out -out report.html
open report.html
```

### Opção 2: Usar como Biblioteca

```bash
go get github.com/rayque.oliveira/coverage-report-generator/pkg/coverage
```

```go
package main

import (
    "os"
    "github.com/rayque.oliveira/coverage-report-generator/pkg/coverage"
)

func main() {
    file, _ := os.Open("coverage.out")
    defer file.Close()

    cov, _ := coverage.ParseCoverageFile(file)
    gen := coverage.NewHTMLGenerator(cov)
    
    output, _ := os.Create("report.html")
    defer output.Close()
    
    gen.Generate(output)
}
```

### Opção 3: Clone e Use

```bash
git clone https://github.com/rayque.oliveira/coverage-report-generator.git
cd coverage-report-generator

# Compilar
go build -o coverage-report ./cmd/coverage-report

# Usar
./coverage-report -in coverage.out -out report.html
```

## 3 Passos Simples

### 1️⃣ Gerar Cobertura

```bash
go test -coverpkg=./... -coverprofile=coverage.out ./...
```

### 2️⃣ Gerar Relatório

```bash
coverage-report
# ou
go run ./cmd/coverage-report
```

### 3️⃣ Abrir no Navegador

```bash
# macOS
open coverage-report.html

# Linux
xdg-open coverage-report.html

# Windows
start coverage-report.html
```

**Pronto! 🎉**

## Opções CLI

```bash
coverage-report -in coverage.out -out report.html

# Flags:
#   -in string     Arquivo de entrada (default "coverage.out")
#   -out string    Arquivo HTML de saída (default "coverage-report.html")
```

## Características

✨ **O que você vai ver:**
- 📊 Dashboard com estatísticas
- 🔍 Busca em tempo real
- 🌳 Árvore de arquivos interativa
- 🎯 Badges coloridas por cobertura
- 📈 Responsivo (funciona em mobile)
- ⚡ Gerado em <300ms

## Próximos Passos

- 📖 Leia o [README](./README.md) para mais detalhes
- 🔧 Consulte [Guia de Integração](./docs/COVERAGE_GUIDE.md) para CI/CD
- 💡 Veja [Features](./docs/FEATURES.md) para conhecer tudo
- 📚 Acesse [API Reference](./pkg/coverage/README.md) para usar como biblioteca

## Suporte

Dúvidas? Consulte a [documentação completa](./docs/)!

---

**Desenvolvido com ❤️ para qualidade**

