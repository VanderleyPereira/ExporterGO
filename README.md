# 🚀 Automacao-Go: Exportador Dinâmico para MySQL

Uma ferramenta CLI performática em Go para extrair dados de bancos MySQL direto para os formatos **CSV, Excel e PDF**.

O diferencial deste projeto é ser **TOTALMENTE DINÂMICO**: ele não utiliza structs fixas para as tabelas. Você aponta para qualquer banco de dados e ele descobre as colunas e tipos de dados em tempo real.

---

## 🛠️ O que ele faz?

- **Descoberta Automática**: Identifica colunas, tipos e registros sem necessidade de criar modelos para cada tabela.
- **Exportação Multiformato**: Gera arquivos `.csv`, `.xlsx` (Excel estilizado) e `.pdf` (com paginação inteligente).
- **Suporte Total a Acentos (UTF-8)**: Implementação robusta para que caracteres especiais como `ã`, `õ`, `é` e o `~` apareçam corretamente no PDF.
- **Tratamento Inteligente de Dados**: Converte automaticamente tipos complexos do MySQL (como `[]byte` em campos VARCHAR/TEXT) para texto legível.
- **Queries Customizadas**: Permite executar comandos SQL `SELECT` manuais com filtros e joins, exportando o resultado final.

---

## 📁 Estrutura do Projeto

```
automacao-go/
├── database/
│   └── connection.go     # Core: Conexão, Listagem e Queries Dinâmicas
├── exporter/
│   ├── interface.go      # Interface comum para os exportadores
│   ├── csv.go            # Geração de CSV com BOM (compatível com Excel)
│   ├── excel.go          # Geração de Planilhas com Estilo (Zebrado e Cabeçalho)
│   └── pdf.go            # Geração de PDF com tradutor Unicode e Rodapé
├── main.go               # Interface de Terminal (CLI) e fluxos de usuário
├── go.mod                # Dependências do projeto
└── go.sum                # Checksum das dependências
```

---

## ⚡ Como Usar

### 1. Preparação

Certifique-se de ter o Go instalado e as dependências baixadas:

```bash
go mod tidy
```

### 2. Execução

Inicie a ferramenta via terminal:

```bash
go run main.go
```

### 3. Fluxo de Operação

1. **Conexão**: Informe Host, Usuário, Senha e o Nome do Banco.
2. **Menu**:
   - `1`: Exportar uma tabela completa (selecionando da lista).
   - `2`: Executar uma Query SQL customizada.
   - `3`: Sair.
3. **Exportação**: Escolha o formato desejado e o arquivo será gerado com um _timestamp_ automático no nome.

---

## 📦 Tecnologias Utilizadas

- **Go (Golang)**: Core do sistema.
- **go-sql-driver/mysql**: Driver de conexão.
- **excelize/v2**: Motor de geração de Excel.
- **gofpdf**: Geração de documentos PDF.

---

## � Exemplos de Queries Customizadas

Você pode rodar queries complexas diretamente no terminal:

```sql
-- Exemplo de Relatório de Vendas com Join
SELECT u.nome, p.descricao, v.valor_total
FROM vendas v
JOIN usuarios u ON u.id = v.usuario_id
JOIN produtos p ON p.id = v.produto_id
WHERE v.data > '2024-01-01'
ORDER BY v.valor_total DESC;
```

---

## � Notas de Versão Final

- Projeto focado em performance via CLI (Interface de Linha de Comando).
- Removidos modelos estáticos para garantir máxima flexibilidade.
- Implementada tradução Unicode para PDFs perfeitos em Português.
