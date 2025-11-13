package coverage

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
)

// CoverageBlock representa um bloco de cobertura no código
type CoverageBlock struct {
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
	NumStmt   int
	Count     int
}

// FileCoverage representa a cobertura de um arquivo
type FileCoverage struct {
	FilePath    string
	FileName    string
	Blocks      []CoverageBlock
	TotalStmt   int
	CoveredStmt int
	Coverage    float64 // percentual de 0-100
}

// ProjectCoverage representa a cobertura de todo o projeto
type ProjectCoverage struct {
	Mode  string
	Files map[string]*FileCoverage
}

// ParseCoverageFile lê e parseia um arquivo de cobertura do Go
func ParseCoverageFile(reader io.Reader) (*ProjectCoverage, error) {
	coverage := &ProjectCoverage{
		Files: make(map[string]*FileCoverage),
	}

	scanner := bufio.NewScanner(reader)

	// Ler primeira linha (modo)
	if !scanner.Scan() {
		return nil, fmt.Errorf("arquivo de cobertura vazio")
	}

	modeLine := scanner.Text()
	if strings.HasPrefix(modeLine, "mode:") {
		coverage.Mode = strings.TrimSpace(strings.TrimPrefix(modeLine, "mode:"))
	}

	// Ler blocos de cobertura
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		block, err := parseCoverageLine(line)
		if err != nil {
			return nil, fmt.Errorf("erro ao parsear linha: %s - %w", line, err)
		}

		// Extrair arquivo do caminho
		filePath := block.FilePath
		if _, ok := coverage.Files[filePath]; !ok {
			coverage.Files[filePath] = &FileCoverage{
				FilePath: filePath,
				FileName: filepath.Base(filePath),
				Blocks:   []CoverageBlock{},
			}
		}

		coverage.Files[filePath].Blocks = append(coverage.Files[filePath].Blocks, block.Block)
		coverage.Files[filePath].TotalStmt += block.Block.NumStmt
		if block.Block.Count > 0 {
			coverage.Files[filePath].CoveredStmt += block.Block.NumStmt
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	// Calcular percentuais
	for _, file := range coverage.Files {
		if file.TotalStmt > 0 {
			file.Coverage = float64(file.CoveredStmt) / float64(file.TotalStmt) * 100
		}
	}

	return coverage, nil
}

type parsedCoverageLine struct {
	FilePath string
	Block    CoverageBlock
}

func parseCoverageLine(line string) (*parsedCoverageLine, error) {
	// Formato: arquivo:startLine.startCol,endLine.endCol numStmt count
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return nil, fmt.Errorf("formato inválido")
	}

	// Parse arquivo e posição
	fileParts := strings.Split(parts[0], ":")
	if len(fileParts) < 2 {
		return nil, fmt.Errorf("formato de arquivo inválido")
	}

	filePath := fileParts[0]
	posStr := fileParts[1]

	// Parse posição: startLine.startCol,endLine.endCol
	posParts := strings.Split(posStr, ",")
	if len(posParts) != 2 {
		return nil, fmt.Errorf("formato de posição inválido")
	}

	startParts := strings.Split(posParts[0], ".")
	endParts := strings.Split(posParts[1], ".")

	if len(startParts) != 2 || len(endParts) != 2 {
		return nil, fmt.Errorf("formato de coordenadas inválido")
	}

	startLine, _ := strconv.Atoi(startParts[0])
	startCol, _ := strconv.Atoi(startParts[1])
	endLine, _ := strconv.Atoi(endParts[0])
	endCol, _ := strconv.Atoi(endParts[1])

	numStmt, _ := strconv.Atoi(parts[1])
	count, _ := strconv.Atoi(parts[2])

	return &parsedCoverageLine{
		FilePath: filePath,
		Block: CoverageBlock{
			StartLine: startLine,
			StartCol:  startCol,
			EndLine:   endLine,
			EndCol:    endCol,
			NumStmt:   numStmt,
			Count:     count,
		},
	}, nil
}

// GetTotalCoverage calcula a cobertura total do projeto
func (pc *ProjectCoverage) GetTotalCoverage() float64 {
	totalStmt := 0
	coveredStmt := 0

	for _, file := range pc.Files {
		totalStmt += file.TotalStmt
		coveredStmt += file.CoveredStmt
	}

	if totalStmt == 0 {
		return 0
	}

	return float64(coveredStmt) / float64(totalStmt) * 100
}

// GetSortedFiles retorna arquivos ordenados por nome
func (pc *ProjectCoverage) GetSortedFiles() []*FileCoverage {
	files := make([]*FileCoverage, 0, len(pc.Files))
	for _, file := range pc.Files {
		files = append(files, file)
	}
	return files
}
