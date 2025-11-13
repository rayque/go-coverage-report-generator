package coverage

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseCoverageFile(t *testing.T) {
	input := `mode: atomic
shipping-management/internal/application/usecases/create_package.go:17.35,22.2 1 2
shipping-management/internal/application/usecases/create_package.go:24.110,27.2 2 2
shipping-management/internal/application/usecases/get_package.go:13.97,17.2 1 2
`

	reader := strings.NewReader(input)
	coverage, err := ParseCoverageFile(reader)

	if err != nil {
		t.Fatalf("erro ao parsear arquivo: %v", err)
	}

	if coverage.Mode != "atomic" {
		t.Errorf("modo esperado 'atomic', obtido '%s'", coverage.Mode)
	}

	if len(coverage.Files) != 2 {
		t.Errorf("esperava 2 arquivos, obteve %d", len(coverage.Files))
	}

	// Verificar arquivo create_package.go
	createPkgFile := coverage.Files["shipping-management/internal/application/usecases/create_package.go"]
	if createPkgFile == nil {
		t.Fatal("arquivo create_package.go não encontrado")
	}

	if len(createPkgFile.Blocks) != 2 {
		t.Errorf("esperava 2 blocos, obteve %d", len(createPkgFile.Blocks))
	}

	if createPkgFile.TotalStmt != 3 {
		t.Errorf("esperava 3 statements, obteve %d", createPkgFile.TotalStmt)
	}

	if createPkgFile.CoveredStmt != 3 {
		t.Errorf("esperava 3 statements cobertos, obteve %d", createPkgFile.CoveredStmt)
	}
}

func TestCoverageCalculation(t *testing.T) {
	input := `mode: atomic
pkg/file.go:1.1,2.2 5 3
pkg/file.go:3.1,4.2 2 0
`

	reader := strings.NewReader(input)
	coverage, err := ParseCoverageFile(reader)

	if err != nil {
		t.Fatalf("erro ao parsear: %v", err)
	}

	file := coverage.Files["pkg/file.go"]

	// 5 de 7 statements cobertos (count > 0 = coberto) = 71.43%
	expectedCov := float64(5) / float64(7) * 100
	if file.Coverage < expectedCov-0.1 || file.Coverage > expectedCov+0.1 {
		t.Errorf("esperava ~%.2f%%, obteve %.2f%%", expectedCov, file.Coverage)
	}
}

func TestHTMLGeneration(t *testing.T) {
	input := `mode: atomic
pkg/file1.go:1.1,2.2 1 1
pkg/file2.go:3.1,4.2 1 0
`

	reader := strings.NewReader(input)
	coverage, err := ParseCoverageFile(reader)

	if err != nil {
		t.Fatalf("erro ao parsear: %v", err)
	}

	generator := NewHTMLGenerator(coverage)
	var buf bytes.Buffer

	err = generator.Generate(&buf)
	if err != nil {
		t.Fatalf("erro ao gerar HTML: %v", err)
	}

	html := buf.String()

	// Verificações básicas
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("HTML não contém DOCTYPE")
	}

	if !strings.Contains(html, "Relatório de Cobertura") {
		t.Error("HTML não contém título esperado")
	}

	if !strings.Contains(html, "file1.go") {
		t.Error("HTML não contém nome do arquivo 1")
	}

	if !strings.Contains(html, "file2.go") {
		t.Error("HTML não contém nome do arquivo 2")
	}

	if !strings.Contains(html, "window.filesData") {
		t.Error("HTML não contém dados dos arquivos")
	}
}

func TestGetTotalCoverage(t *testing.T) {
	input := `mode: atomic
pkg/file1.go:1.1,2.2 10 8
pkg/file2.go:3.1,4.2 10 4
`

	reader := strings.NewReader(input)
	coverage, err := ParseCoverageFile(reader)

	if err != nil {
		t.Fatalf("erro ao parsear: %v", err)
	}

	total := coverage.GetTotalCoverage()
	// 10 cobertos + 10 cobertos = 20 de 20 = 100%
	expected := 100.0
	if total < expected-0.1 || total > expected+0.1 {
		t.Errorf("esperava ~%.1f%%, obteve %.1f%%", expected, total)
	}
}

func TestParseCoverageLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantErr   bool
		wantStart int
		wantEnd   int
	}{
		{
			name:      "linha válida",
			line:      "pkg/file.go:10.5,15.10 2 1",
			wantErr:   false,
			wantStart: 10,
			wantEnd:   15,
		},
		{
			name:    "formato inválido",
			line:    "pkg/file.go 2 1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseCoverageLine(tt.line)

			if (err != nil) != tt.wantErr {
				t.Errorf("erro = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result.Block.StartLine != tt.wantStart {
					t.Errorf("startLine = %d, want %d", result.Block.StartLine, tt.wantStart)
				}
				if result.Block.EndLine != tt.wantEnd {
					t.Errorf("endLine = %d, want %d", result.Block.EndLine, tt.wantEnd)
				}
			}
		})
	}
}

func BenchmarkParseCoverageFile(b *testing.B) {
	input := `mode: atomic
pkg/file1.go:1.1,2.2 1 1
pkg/file2.go:3.1,4.2 1 0
pkg/file3.go:5.1,6.2 1 1
pkg/file4.go:7.1,8.2 1 0
pkg/file5.go:9.1,10.2 1 1
`

	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(input)
		ParseCoverageFile(reader)
	}
}
