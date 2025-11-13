package coverage

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestIntegration testa o fluxo completo
func TestIntegration(t *testing.T) {
	// Criar arquivo temporário
	tmpFile, _ := os.CreateTemp("", "coverage*.out")
	defer os.Remove(tmpFile.Name())

	// Escrever dados de cobertura
	content := `mode: atomic
pkg/file1.go:1.1,2.2 5 5
pkg/file1.go:3.1,4.2 5 3
pkg/file2.go:5.1,6.2 5 0
`
	tmpFile.WriteString(content)
	tmpFile.Close()

	// Abrir e parsear
	file, _ := os.Open(tmpFile.Name())
	cov, err := ParseCoverageFile(file)
	file.Close()

	if err != nil {
		t.Fatalf("erro ao parsear: %v", err)
	}

	// Verificar estatísticas
	if len(cov.Files) != 2 {
		t.Errorf("esperava 2 arquivos, obteve %d", len(cov.Files))
	}

	total := cov.GetTotalCoverage()
	// file1: 5 cobertos + 5 cobertos = 10 de 10 = 100%
	// file2: 0 cobertos = 0 de 5 = 0%
	// total: 10 de 15 = 66.67%
	if total < 66 || total > 67 {
		t.Errorf("cobertura esperada ~66.67%%, obteve %.2f%%", total)
	}

	// Gerar HTML
	generator := NewHTMLGenerator(cov)
	var buf bytes.Buffer
	err = generator.Generate(&buf)

	if err != nil {
		t.Fatalf("erro ao gerar HTML: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("HTML inválido")
	}
}
