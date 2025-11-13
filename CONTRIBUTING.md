# Coverage Report Generator

Repositório da biblioteca Go para geração de relatórios de cobertura HTML interativos.

## 📁 Estrutura

```
coverage-report-generator/
├── pkg/coverage/           # Biblioteca principal
│   ├── parser.go
│   ├── html_generator.go
│   ├── parser_test.go
│   ├── examples_test.go
│   ├── README.md           # API Reference
│   └── LIBRARY.md
├── cmd/coverage-report/    # CLI
│   └── main.go
├── scripts/                # Scripts auxiliares
│   └── generate-coverage-report.sh
├── docs/                   # Documentação
│   ├── COVERAGE_GUIDE.md
│   └── FEATURES.md
├── go.mod
├── go.sum
├── Makefile
├── README.md               # Este arquivo
└── .gitignore
```

## 🚀 Começar Rápido

```bash
# 1. Clonar repositório
git clone https://github.com/rayque.oliveira/coverage-report-generator.git
cd coverage-report-generator

# 2. Instalar
go install ./cmd/coverage-report

# 3. Usar
go test -coverpkg=./... -coverprofile=coverage.out ./...
coverage-report -in coverage.out -out report.html
open report.html
```

## 📚 Documentação

- [README](./README.md) - Visão geral
- [API Reference](./pkg/coverage/README.md) - Documentação completa
- [Guia de Integração](./docs/COVERAGE_GUIDE.md) - CI/CD e integração
- [Features](./docs/FEATURES.md) - Lista de recursos

## 🧪 Testes

```bash
make test          # Rodar testes
make test-coverage # Gerar cobertura
make example       # Gerar exemplo
```

## 📦 Publicação no GitHub

Este repositório está pronto para ser publicado no GitHub. Próximos passos:

1. Criar repositório no GitHub
2. `git init` neste diretório
3. `git remote add origin https://github.com/seu-usuario/coverage-report-generator.git`
4. `git add .`
5. `git commit -m "Initial commit"`
6. `git push -u origin main`

## 📝 Licença

MIT

## 🤝 Contribuir

Contribute com pull requests!

---

Desenvolvido com ❤️ para qualidade

