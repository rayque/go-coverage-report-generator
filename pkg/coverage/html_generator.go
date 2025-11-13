package coverage

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
)

// HTMLGenerator gera relatórios HTML da cobertura
type HTMLGenerator struct {
	coverage *ProjectCoverage
}

// NewHTMLGenerator cria um novo gerador de HTML
func NewHTMLGenerator(coverage *ProjectCoverage) *HTMLGenerator {
	return &HTMLGenerator{coverage: coverage}
}

// Generate gera o HTML e escreve no writer
func (hg *HTMLGenerator) Generate(writer io.Writer) error {
	// Escrever cabeçalho HTML
	if err := hg.writeHeader(writer); err != nil {
		return err
	}

	// Escrever CSS
	if err := hg.writeStyles(writer); err != nil {
		return err
	}

	// Escrever conteúdo
	if err := hg.writeContent(writer); err != nil {
		return err
	}

	// Escrever JavaScript
	if err := hg.writeScripts(writer); err != nil {
		return err
	}

	// Escrever rodapé
	if err := hg.writeFooter(writer); err != nil {
		return err
	}

	return nil
}

func (hg *HTMLGenerator) writeHeader(w io.Writer) error {
	html := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Relatório de Cobertura de Testes</title>
</head>
<body>
`
	_, err := io.WriteString(w, html)
	return err
}

func (hg *HTMLGenerator) writeStyles(w io.Writer) error {
	css := `<style>
    * {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
    }

    body {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
        background: #f6f8fa;
        color: #24292e;
        line-height: 1.5;
    }

    .container {
        max-width: 1280px;
        margin: 0 auto;
        padding: 20px;
    }

    header {
        background: white;
        border-bottom: 1px solid #e1e4e8;
        margin-bottom: 20px;
        padding: 20px 0;
        box-shadow: 0 1px 3px rgba(0,0,0,0.05);
    }

    .header-content {
        padding: 0 20px;
    }

    h1 {
        font-size: 24px;
        font-weight: 600;
        margin-bottom: 10px;
    }

    .coverage-badge {
        display: inline-block;
        padding: 8px 16px;
        border-radius: 6px;
        font-weight: 600;
        font-size: 14px;
        margin-top: 10px;
    }

    .coverage-excellent {
        background: #dcffe4;
        color: #0d643d;
        border: 1px solid #34d399;
    }

    .coverage-good {
        background: #cce5ff;
        color: #0366d6;
        border: 1px solid #0366d6;
    }

    .coverage-fair {
        background: #fff8c5;
        color: #856404;
        border: 1px solid #ffc107;
    }

    .coverage-poor {
        background: #ffeef0;
        color: #6f42c1;
        border: 1px solid #ff6a88;
    }

    .controls {
        display: flex;
        gap: 10px;
        margin-top: 15px;
        align-items: center;
    }

    .search-box {
        flex: 1;
        max-width: 300px;
    }

    input[type="text"] {
        width: 100%;
        padding: 8px 12px;
        border: 1px solid #e1e4e8;
        border-radius: 6px;
        font-size: 14px;
        transition: border-color 0.2s;
    }

    input[type="text"]:focus {
        outline: none;
        border-color: #0366d6;
        box-shadow: 0 0 0 3px rgba(3, 102, 214, 0.1);
    }

    button {
        padding: 8px 16px;
        background: #28a745;
        color: white;
        border: none;
        border-radius: 6px;
        cursor: pointer;
        font-size: 14px;
        font-weight: 500;
        transition: background 0.2s;
    }

    button:hover {
        background: #218838;
    }

    .main-content {
        display: flex;
        gap: 20px;
        margin-top: 20px;
    }

    .file-tree {
        flex: 0 0 300px;
        background: white;
        border: 1px solid #e1e4e8;
        border-radius: 6px;
        overflow-y: auto;
        max-height: calc(100vh - 200px);
        position: sticky;
        top: 20px;
    }

    .file-list {
        list-style: none;
    }

    .file-item {
        border-bottom: 1px solid #e1e4e8;
        transition: background 0.1s;
    }

    .file-item:hover {
        background: #f6f8fa;
    }

    .file-item.active {
        background: #e1ecf7;
        border-left: 4px solid #0366d6;
    }

    .file-link {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 12px;
        cursor: pointer;
        text-decoration: none;
        color: #24292e;
        font-size: 13px;
        user-select: none;
    }

    .file-name {
        display: flex;
        align-items: center;
        flex: 1;
        overflow: hidden;
    }

    .file-name-text {
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .file-coverage-badge {
        padding: 2px 8px;
        border-radius: 4px;
        font-size: 12px;
        font-weight: 600;
        margin-left: 8px;
        flex-shrink: 0;
    }

    .main-panel {
        flex: 1;
        background: white;
        border: 1px solid #e1e4e8;
        border-radius: 6px;
        overflow: hidden;
        display: flex;
        flex-direction: column;
        min-height: 500px;
    }

    .file-header {
        padding: 16px;
        border-bottom: 1px solid #e1e4e8;
        background: #f6f8fa;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .file-header-info {
        display: flex;
        flex-direction: column;
        gap: 5px;
    }

    .file-header-title {
        font-size: 14px;
        font-weight: 600;
    }

    .file-path {
        font-size: 12px;
        color: #666;
        font-family: monospace;
    }

    .code-view {
        flex: 1;
        overflow: auto;
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
        font-size: 12px;
        line-height: 1.6;
    }

    .code-line {
        display: flex;
        border-bottom: 1px solid #eee;
        transition: background 0.1s;
    }

    .code-line:hover {
        background: #f5f5f5;
    }

    .line-number {
        flex: 0 0 50px;
        padding: 2px 10px;
        text-align: right;
        background: #f6f8fa;
        color: #666;
        user-select: none;
        border-right: 1px solid #e1e4e8;
    }

    .coverage-indicator {
        flex: 0 0 8px;
        background: #eee;
        transition: background 0.1s;
    }

    .covered .coverage-indicator {
        background: #dcffe4;
    }

    .uncovered .coverage-indicator {
        background: #ffeef0;
    }

    .code-content {
        flex: 1;
        padding: 2px 16px;
        white-space: pre-wrap;
        word-break: break-word;
        color: #24292e;
    }

    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 400px;
        color: #666;
    }

    .empty-state-icon {
        font-size: 48px;
        margin-bottom: 16px;
    }

    .stats-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 16px;
        margin-bottom: 20px;
    }

    .stat-card {
        background: white;
        border: 1px solid #e1e4e8;
        border-radius: 6px;
        padding: 16px;
    }

    .stat-label {
        font-size: 12px;
        color: #666;
        text-transform: uppercase;
        letter-spacing: 0.5px;
        font-weight: 600;
        margin-bottom: 8px;
    }

    .stat-value {
        font-size: 28px;
        font-weight: 600;
        color: #0366d6;
    }

    .stat-subtext {
        font-size: 12px;
        color: #666;
        margin-top: 4px;
    }

    .progress-bar {
        width: 100%;
        height: 8px;
        background: #e1e4e8;
        border-radius: 4px;
        overflow: hidden;
        margin-top: 8px;
    }

    .progress-fill {
        height: 100%;
        background: #28a745;
        transition: width 0.3s;
    }

    .sort-controls {
        display: flex;
        gap: 10px;
        align-items: center;
    }

    .sort-btn {
        padding: 6px 12px;
        background: #f6f8fa;
        color: #24292e;
        border: 1px solid #e1e4e8;
        border-radius: 6px;
        cursor: pointer;
        font-size: 13px;
        font-weight: 500;
        transition: background 0.2s;
    }

    .sort-btn:hover,
    .sort-btn.active {
        background: #0366d6;
        color: white;
        border-color: #0366d6;
    }

    @media (max-width: 768px) {
        .main-content {
            flex-direction: column;
        }

        .file-tree {
            max-height: 300px;
            position: static;
            flex: 0 0 auto;
        }

        .stats-grid {
            grid-template-columns: repeat(2, 1fr);
        }
    }
</style>
`
	_, err := io.WriteString(w, css)
	return err
}

func (hg *HTMLGenerator) writeContent(w io.Writer) error {
	totalCoverage := hg.coverage.GetTotalCoverage()
	files := hg.coverage.Files

	// Contar estatísticas
	totalFiles := len(files)
	totalStmt := 0
	coveredStmt := 0

	for _, file := range files {
		totalStmt += file.TotalStmt
		coveredStmt += file.CoveredStmt
	}

	// Determinar cor da cobertura
	coverageClass := getCoverageClass(totalCoverage)
	coverageText := fmt.Sprintf("%.1f%%", totalCoverage)

	html := fmt.Sprintf(`<header>
    <div class="header-content">
        <h1>📊 Relatório de Cobertura de Testes</h1>
        <div class="coverage-badge %s">Cobertura Total: %s</div>
        <div class="controls">
            <div class="search-box">
                <input type="text" id="searchInput" placeholder="Buscar arquivo...">
            </div>
            <div class="sort-controls">
                <button class="sort-btn active" onclick="sortFiles('name')">Nome</button>
                <button class="sort-btn" onclick="sortFiles('coverage')">Cobertura</button>
            </div>
        </div>
    </div>
</header>

<div class="container">
    <div class="stats-grid">
        <div class="stat-card">
            <div class="stat-label">Total de Arquivos</div>
            <div class="stat-value">%d</div>
        </div>
        <div class="stat-card">
            <div class="stat-label">Linhas Cobertas</div>
            <div class="stat-value">%d</div>
            <div class="stat-subtext">de %d</div>
            <div class="progress-bar">
                <div class="progress-fill" style="width: %.1f%%"></div>
            </div>
        </div>
        <div class="stat-card">
            <div class="stat-label">Modo</div>
            <div class="stat-value">%s</div>
        </div>
    </div>

    <div class="main-content">
        <div class="file-tree">
            <ul class="file-list" id="fileList">
`, coverageClass, coverageText, totalFiles, coveredStmt, totalStmt,
		(float64(coveredStmt)/float64(totalStmt))*100, hg.coverage.Mode)

	if _, err := io.WriteString(w, html); err != nil {
		return err
	}

	// Escrever lista de arquivos
	fileList := make([]*FileCoverage, 0, len(files))
	for _, f := range files {
		fileList = append(fileList, f)
	}

	sort.Slice(fileList, func(i, j int) bool {
		return fileList[i].FilePath < fileList[j].FilePath
	})

	for i, file := range fileList {
		fcClass := getCoverageClass(file.Coverage)
		covText := fmt.Sprintf("%.0f%%", file.Coverage)

		active := ""
		if i == 0 {
			active = " active"
		}

		fileHTML := fmt.Sprintf(`                <li class="file-item%s">
                    <a href="#" class="file-link" onclick="loadFile('%s', '%s', '%d', '%d'); return false;">
                        <span class="file-name">
                            <span class="file-name-text">📄 %s</span>
                            <span class="file-coverage-badge %s">%s</span>
                        </span>
                    </a>
                </li>
`, active, strings.ReplaceAll(file.FilePath, "'", "\\'"),
			file.FilePath, file.CoveredStmt, file.TotalStmt,
			file.FileName, fcClass, covText)

		if _, err := io.WriteString(w, fileHTML); err != nil {
			return err
		}
	}

	mainHTML := `            </ul>
        </div>

        <div class="main-panel">
            <div id="content">
                <div class="empty-state">
                    <div class="empty-state-icon">👈</div>
                    <p>Selecione um arquivo para visualizar a cobertura</p>
                </div>
            </div>
        </div>
    </div>
</div>
`

	_, err := io.WriteString(w, mainHTML)
	return err
}

func (hg *HTMLGenerator) writeScripts(w io.Writer) error {
	// Preparar dados dos arquivos
	filesData := "window.filesData = {\n"

	fileList := make([]*FileCoverage, 0, len(hg.coverage.Files))
	for _, f := range hg.coverage.Files {
		fileList = append(fileList, f)
	}

	sort.Slice(fileList, func(i, j int) bool {
		return fileList[i].FilePath < fileList[j].FilePath
	})

	for i, file := range fileList {
		blocksMap := make(map[int]int) // line -> count
		for _, block := range file.Blocks {
			for line := block.StartLine; line <= block.EndLine; line++ {
				if block.Count > 0 {
					blocksMap[line] = 1
				} else {
					if _, exists := blocksMap[line]; !exists {
						blocksMap[line] = 0
					}
				}
			}
		}

		// Serializar blocos
		var blockStrs []string
		for line, count := range blocksMap {
			blockStrs = append(blockStrs, fmt.Sprintf("%d:%d", line, count))
		}

		filesData += fmt.Sprintf("    '%s': {\n        filePath: '%s',\n        fileName: '%s',\n        coverage: %.2f,\n        covered: %d,\n        total: %d,\n        blocks: '%s'\n    }",
			strings.ReplaceAll(file.FilePath, "'", "\\'"),
			strings.ReplaceAll(file.FilePath, "'", "\\'"),
			strings.ReplaceAll(file.FileName, "'", "\\'"),
			file.Coverage,
			file.CoveredStmt,
			file.TotalStmt,
			strings.Join(blockStrs, ","),
		)

		if i < len(fileList)-1 {
			filesData += ",\n"
		}
	}
	filesData += "\n};\n"

	if _, err := io.WriteString(w, filesData); err != nil {
		return err
	}

	js := `<script>
let currentFile = null;
let sortBy = 'name';

function sortFiles(by) {
    sortBy = by;
    const fileList = document.getElementById('fileList');
    const files = Array.from(fileList.children);

    files.sort((a, b) => {
        if (by === 'name') {
            return a.textContent.localeCompare(b.textContent);
        } else if (by === 'coverage') {
            const aCov = parseFloat(a.querySelector('.file-coverage-badge').textContent);
            const bCov = parseFloat(b.querySelector('.file-coverage-badge').textContent);
            return bCov - aCov;
        }
    });

    files.forEach(f => fileList.appendChild(f));

    // Atualizar botões
    document.querySelectorAll('.sort-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    event.target.classList.add('active');
}

function loadFile(filePath, fileName, covered, total) {
    currentFile = filePath;

    // Atualizar seleção
    document.querySelectorAll('.file-item').forEach(item => {
        item.classList.remove('active');
    });
    event.target.closest('.file-item').classList.add('active');

    // Preparar header
    const coverage = total > 0 ? ((covered / total) * 100).toFixed(1) : 0;
    const headerHTML = '<div class="file-header">' +
        '<div class="file-header-info">' +
        '<div class="file-header-title">' + fileName + '</div>' +
        '<div class="file-path">' + filePath + '</div>' +
        '</div>' +
        '<div>' +
        '<span style="font-weight: 600; color: #0366d6;">' + coverage + '% (' + covered + '/' + total + ')</span>' +
        '</div>' +
        '</div>' +
        '<div class="code-view" id="codeView">' +
        '</div>';

    const content = document.getElementById('content');
    content.innerHTML = headerHTML;

    // Renderizar código (simulado - em produção teria que ler o arquivo)
    const codeView = document.getElementById('codeView');
    const blockData = window.filesData[filePath];
    
    if (blockData) {
        const blocks = {};
        blockData.blocks.split(',').forEach(b => {
            const [line, count] = b.split(':').map(Number);
            blocks[line] = count;
        });

        // Simular conteúdo do arquivo
        for (let i = 1; i <= total; i++) {
            const isActive = blocks[i] !== undefined;
            const isCovered = blocks[i] === 1;

            const lineClass = isActive ? (isCovered ? 'covered' : 'uncovered') : '';
            const lineHTML = '<div class="code-line ' + lineClass + '">' +
                '<div class="line-number">' + i + '</div>' +
                '<div class="coverage-indicator"></div>' +
                '<div class="code-content">// Linha ' + i + '</div>' +
                '</div>';
            codeView.innerHTML += lineHTML;
        }
    }
}

document.getElementById('searchInput').addEventListener('input', function(e) {
    const query = e.target.value.toLowerCase();
    const items = document.querySelectorAll('.file-item');

    items.forEach(item => {
        const fileName = item.querySelector('.file-name-text').textContent.toLowerCase();
        item.style.display = fileName.includes(query) ? '' : 'none';
    });
});

// Carregar primeiro arquivo ao iniciar
document.addEventListener('DOMContentLoaded', function() {
    const firstFile = document.querySelector('.file-item.active .file-link');
    if (firstFile) {
        firstFile.click();
    }
});
</script>
`
	_, err := io.WriteString(w, js)
	return err
}

func (hg *HTMLGenerator) writeFooter(w io.Writer) error {
	footer := `    </body>
</html>
`
	_, err := io.WriteString(w, footer)
	return err
}

func getCoverageClass(coverage float64) string {
	switch {
	case coverage >= 80:
		return "coverage-excellent"
	case coverage >= 60:
		return "coverage-good"
	case coverage >= 40:
		return "coverage-fair"
	default:
		return "coverage-poor"
	}
}

func roundFloat(f float64, decimals int) float64 {
	shift := math.Pow(10, float64(decimals))
	return math.Round(f*shift) / shift
}
