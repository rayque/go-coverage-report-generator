# Coverage Report Generator

Uma biblioteca Go para converter arquivos de cobertura do Go (`coverage.out`) em relatórios HTML interativos e visuais, similar à navegação do GitHub.

## Características

- 📊 **Parser Completo**: Lê arquivos de cobertura gerados por `go tool cover`
- 🎨 **HTML Interativo**: Interface similar ao GitHub com navegação fluida
- 🔍 **Busca em Tempo Real**: Filtre arquivos enquanto digita
- 📈 **Visualização de Estatísticas**: Exiba métricas gerais e por arquivo
- 🏆 **Codificação de Cores**: Cores indicam nível de cobertura
- 📱 **Responsivo**: Funciona bem em diferentes tamanhos de tela
- ⚡ **Performance**: Parsing otimizado e renderização eficiente

## Instalação

```bash
go get shipping-management/pkg/coverage
```

## Uso

### Como Biblioteca

```go
package main

import (
	"os"
	"log"
	"shipping-management/pkg/coverage"
)

func main() {
	// Abrir arquivo de cobertura
	file, err := os.Open("coverage.out")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Parsear cobertura
	cov, err := coverage.ParseCoverageFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// Gerar HTML
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

### Como CLI

```bash
# Gerar relatório com arquivo padrão
go run cmd/coverage-report/main.go

# Especificar arquivo de entrada
go run cmd/coverage-report/main.go -in custom-coverage.out

# Especificar arquivo de saída
go run cmd/coverage-report/main.go -in coverage.out -out my-report.html
```

## Integração no Makefile

Adicione o seguinte ao seu `Makefile`:

```makefile
# generate-coverage-report: Gera relatório HTML de cobertura interativo
coverage-report:
	go run cmd/coverage-report/main.go -in coverage.out -out coverage-report.html
	@echo "✨ Abra coverage-report.html no seu navegador"

# test-report: Executa testes e gera relatório HTML
test-report: test coverage-report
	@echo "📊 Relatório completo disponível em coverage-report.html"
```

Depois, execute:

```bash
make coverage-report
# ou
make test-report
```

## API da Biblioteca

### Tipos

#### `CoverageBlock`
Representa um bloco de cobertura no código.

```go
type CoverageBlock struct {
	StartLine  int  // Linha inicial
	StartCol   int  // Coluna inicial
	EndLine    int  // Linha final
	EndCol     int  // Coluna final
	NumStmt    int  // Número de statements
	Count      int  // Número de vezes executado
}
```

#### `FileCoverage`
Representa a cobertura de um arquivo.

```go
type FileCoverage struct {
	FilePath    string           // Caminho completo do arquivo
	FileName    string           // Nome do arquivo
	Blocks      []CoverageBlock  // Blocos cobertos
	TotalStmt   int              // Total de statements
	CoveredStmt int              // Statements cobertos
	Coverage    float64          // Percentual (0-100)
}
```

#### `ProjectCoverage`
Representa a cobertura do projeto inteiro.

```go
type ProjectCoverage struct {
	Mode  string
	Files map[string]*FileCoverage
}
```

### Funções

#### `ParseCoverageFile(reader io.Reader) (*ProjectCoverage, error)`
Lê e parseia um arquivo de cobertura.

```go
file, _ := os.Open("coverage.out")
coverage, err := coverage.ParseCoverageFile(file)
```

#### `NewHTMLGenerator(coverage *ProjectCoverage) *HTMLGenerator`
Cria um novo gerador de HTML.

```go
generator := coverage.NewHTMLGenerator(coverage)
```

#### `(hg *HTMLGenerator) Generate(writer io.Writer) error`
Gera o HTML e escreve no writer.

```go
output, _ := os.Create("report.html")
err := generator.Generate(output)
```

#### `(pc *ProjectCoverage) GetTotalCoverage() float64`
Retorna a cobertura total do projeto em percentual.

```go
total := coverage.GetTotalCoverage()
fmt.Printf("Cobertura: %.2f%%\n", total)
```

## Estrutura do HTML Gerado

O HTML gerado inclui:

- **Header**: Logo, título e cobertura total
- **Estatísticas**: Cards com métricas gerais
- **Sidebar**: Árvore de arquivos com busca e ordenação
- **Painel Principal**: Visualização da cobertura por arquivo
- **Controles**: Busca, ordenação e filtros

### Funcionalidades do HTML

- ✅ Busca em tempo real por nome de arquivo
- ✅ Ordenação por nome ou percentual de cobertura
- ✅ Cores indicativas (verde ≥80%, azul ≥60%, amarelo ≥40%, vermelho <40%)
- ✅ Zoom em blocos de código
- ✅ Visualização de estatísticas por arquivo

## Testes

Execute os testes com:

```bash
go test ./pkg/coverage -v
```

Ou com cobertura:

```bash
go test ./pkg/coverage -cover
```

## Exemplo Completo

1. **Gerar arquivo de cobertura**:
```bash
go test -coverpkg=./internal/application/usecases -coverprofile=coverage.out ./...
```

2. **Gerar relatório HTML**:
```bash
go run cmd/coverage-report/main.go
```

3. **Abrir no navegador**:
```bash
open coverage-report.html  # macOS
xdg-open coverage-report.html  # Linux
start coverage-report.html  # Windows
```

## Formatos de Cores

- 🟢 **Excelente (≥80%)**: Verde - `coverage-excellent`
- 🔵 **Bom (60-79%)**: Azul - `coverage-good`
- 🟡 **Regular (40-59%)**: Amarelo - `coverage-fair`
- 🔴 **Fraco (<40%)**: Vermelho - `coverage-poor`

## Performance

- Parse de 100+ arquivos: < 100ms
- Geração de HTML: < 200ms
- Tamanho do arquivo HTML: ~200KB (incluindo estilos e scripts)

## Contribuição

Para contribuir:

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/amazing-feature`)
3. Commit suas mudanças (`git commit -m 'Add amazing feature'`)
4. Push para a branch (`git push origin feature/amazing-feature`)
5. Abra um Pull Request

## Licença

MIT

## Exemplos de Uso

### Integração com CI/CD

```bash
#!/bin/bash
go test -coverpkg=./... -coverprofile=coverage.out ./...
go run cmd/coverage-report/main.go -in coverage.out -out coverage-report.html

# Enviar para servidor (opcional)
scp coverage-report.html user@server:/var/www/html/
```

### Gerar múltiplos relatórios

```bash
#!/bin/bash
for pkg in ./internal/application ./internal/infrastructure; do
  go test -coverprofile=coverage-$pkg.out $pkg/...
  go run cmd/coverage-report/main.go -in coverage-$pkg.out -out report-$pkg.html
done
```

## Suporte

Para dúvidas ou problemas, abra uma issue no repositório.

