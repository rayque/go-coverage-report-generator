#!/bin/bash

# Script para gerar relatório de cobertura HTML interativo
# Uso: ./generate-coverage-report.sh

set -e

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}📊 Script de Geração de Relatório de Cobertura${NC}\n"

# Função para exibir uso
usage() {
    echo "Uso: $0 [opções]"
    echo ""
    echo "Opções:"
    echo "  -p, --package    Pacotes a cobrir (padrão: ./internal/...)"
    echo "  -o, --output     Arquivo HTML de saída (padrão: coverage-report.html)"
    echo "  -c, --coverage   Arquivo de cobertura (padrão: coverage.out)"
    echo "  -h, --help       Exibir esta mensagem"
    echo ""
    echo "Exemplos:"
    echo "  $0"
    echo "  $0 --package ./internal/application --output report.html"
    exit 0
}

# Valores padrão
PACKAGE="./..."
OUTPUT="coverage-report.html"
COVERAGE="coverage.out"

# Parse de argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--package)
            PACKAGE="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT="$2"
            shift 2
            ;;
        -c|--coverage)
            COVERAGE="$2"
            shift 2
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo "Opção desconhecida: $1"
            usage
            ;;
    esac
done

# Passo 1: Gerar cobertura
echo -e "${YELLOW}1️⃣  Executando testes e gerando cobertura...${NC}"
go test -coverpkg="$PACKAGE" -coverprofile="$COVERAGE" ./...

if [ ! -f "$COVERAGE" ]; then
    echo -e "${RED}❌ Erro: arquivo de cobertura não foi gerado${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Cobertura gerada: $COVERAGE${NC}\n"

# Passo 2: Exibir resumo da cobertura
echo -e "${YELLOW}2️⃣  Resumo de cobertura:${NC}"
go tool cover -func="$COVERAGE" | tail -1
echo ""

# Passo 3: Gerar relatório HTML
echo -e "${YELLOW}3️⃣  Gerando relatório HTML interativo...${NC}"
go run cmd/coverage-report/main.go -in "$COVERAGE" -out "$OUTPUT"

echo ""
echo -e "${GREEN}✨ Relatório gerado com sucesso!${NC}"
echo -e "${BLUE}📂 Arquivo: $OUTPUT${NC}"
echo ""
echo "Para abrir no navegador:"
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo -e "  ${YELLOW}open $OUTPUT${NC}"
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo -e "  ${YELLOW}xdg-open $OUTPUT${NC}"
elif [[ "$OSTYPE" == "msys" ]]; then
    echo -e "  ${YELLOW}start $OUTPUT${NC}"
fi

