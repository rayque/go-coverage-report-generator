# 📋 COPY & PASTE - Comandos Prontos

## 🔗 Comandos para Publicar no GitHub

### Passo 1: Iniciar Git Localmente

```bash
cd /home/rayque.oliveira/projects/coverage-report-generator
git init
```

### Passo 2: Adicionar Todos os Arquivos

```bash
git add .
```

### Passo 3: Fazer Primeiro Commit

```bash
git commit -m "Initial commit: Coverage Report Generator - Professional Go library for generating interactive HTML coverage reports"
```

### Passo 4: Adicionar Remote (SUBSTITUA SEU_USUARIO)

```bash
git remote add origin https://github.com/SEU_USUARIO/coverage-report-generator.git
```

### Passo 5: Renomear Branch e Fazer Push

```bash
git branch -M main
git push -u origin main
```

### Passo 6 (Opcional): Criar Release

```bash
git tag -a v1.0.0 -m "First stable release"
git push origin v1.0.0
```

---

## ✅ Verificações Antes de Publicar

### Testar Compilação

```bash
cd /home/rayque.oliveira/projects/coverage-report-generator
go test ./pkg/coverage -v
```

### Compilar CLI

```bash
go build -o coverage-report ./cmd/coverage-report
```

### Verificar Arquivos

```bash
git status
```

### Ver Log

```bash
git log --oneline
```

---

## 🔧 Comandos Úteis Após Publicar

### Instalar CLI Globalmente

```bash
go install github.com/SEU_USUARIO/coverage-report-generator/cmd/coverage-report@latest
```

### Usar Como Biblioteca

```bash
go get github.com/SEU_USUARIO/coverage-report-generator/pkg/coverage
```

---

## 📝 IMPORTANTE: Editar go.mod

Se estiver usando um usuário GitHub DIFERENTE de `rayque.oliveira`, edite o `go.mod`:

```bash
# Abra o arquivo
nano /home/rayque.oliveira/projects/coverage-report-generator/go.mod

# Altere a primeira linha de:
module github.com/rayque.oliveira/coverage-report-generator

# Para:
module github.com/SEU_USUARIO/coverage-report-generator
```

---

## 🎯 Sequência Completa (Copy & Paste)

```bash
# 1. Navegar
cd /home/rayque.oliveira/projects/coverage-report-generator

# 2. Inicializar git
git init

# 3. Adicionar arquivos
git add .

# 4. Primeiro commit
git commit -m "Initial commit: Coverage Report Generator"

# 5. Conectar ao GitHub (ALTERE SEU_USUARIO)
git remote add origin https://github.com/SEU_USUARIO/coverage-report-generator.git

# 6. Renomear branch
git branch -M main

# 7. Fazer push
git push -u origin main

# 8 (Opcional). Criar tag
git tag -a v1.0.0 -m "First stable release"
git push origin v1.0.0
```

---

## 📋 Checklist Antes de Fazer Push

- [ ] Criar repositório no GitHub em https://github.com/new
- [ ] Verificar que go.mod tem o módulo correto
- [ ] Rodar `go test ./pkg/coverage -v` (todos devem passar)
- [ ] Verificar que os arquivos estão em `git status`
- [ ] Alterar `SEU_USUARIO` nos comandos do GitHub

---

## 🚀 Resumo

```
1. CRIAR no GitHub
2. cd /home/rayque.oliveira/projects/coverage-report-generator
3. git init && git add . && git commit -m "Initial commit"
4. git remote add origin https://github.com/SEU_USUARIO/coverage-report-generator.git
5. git branch -M main && git push -u origin main
6. ✅ PRONTO!
```

---

**Copie, cole e execute! 🎉**

