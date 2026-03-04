package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"automacao-go/database"
	"automacao-go/exporter"
)

var reader = bufio.NewReader(os.Stdin)

func input(prompt string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func main() {
	printBanner()

	fmt.Println("┌─────────────────────────────────┐")
	fmt.Println("│       CONFIGURAR CONEXÃO        │")
	fmt.Println("└─────────────────────────────────┘")

	cfg := database.Config{
		Host:     inputDefault("  Host     [localhost]: ", "localhost"),
		Port:     inputDefault("  Porta    [3306]: ", "3306"),
		User:     inputDefault("  Usuário  [root]: ", "root"),
		Password: input("  Senha: "),
		DBName:   input("  Banco de dados: "),
	}

	if err := database.Connect(cfg); err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	for {
		fmt.Println()
		fmt.Println("┌─────────────────────────────────┐")
		fmt.Println("│           MENU PRINCIPAL        │")
		fmt.Println("├─────────────────────────────────┤")
		fmt.Println("│  1. Exportar tabela completa    │")
		fmt.Println("│  2. Query customizada           │")
		fmt.Println("│  3. Sair                        │")
		fmt.Println("└─────────────────────────────────┘")

		op := input("Escolha: ")

		switch op {
		case "1":
			exportTable()
		case "2":
			exportCustomQuery()
		case "3":
			fmt.Println("👋 Até logo!")
			os.Exit(0)
		default:
			fmt.Println("❌ Opção inválida.")
		}
	}
}

// lista todas as tabelas do banco
func listTables() {
	tables, err := database.ListTables()
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	fmt.Printf("\n📋 Tabelas em '%s':\n", "banco")
	for i, t := range tables {
		fmt.Printf("  %d. %s\n", i+1, t)
	}
}

// exportTable exporta uma tabela completa
func exportTable() {
	tables, err := database.ListTables()
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	fmt.Println("\n📋 Tabelas disponíveis:")
	for i, t := range tables {
		fmt.Printf("  %d. %s\n", i+1, t)
	}

	choice := input("\nNúmero ou nome da tabela: ")

	var tableName string
	if idx, err := strconv.Atoi(choice); err == nil {
		if idx < 1 || idx > len(tables) {
			fmt.Println("❌ Número inválido.")
			return
		}
		tableName = tables[idx-1]
	} else {
		tableName = choice
	}

	fmt.Printf("🔍 Buscando dados de '%s'...\n", tableName)
	headers, rows, err := database.QueryTable(tableName)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	fmt.Printf("📊 %d colunas | %d registros encontrados\n", len(headers), len(rows))
	doExport(headers, rows, tableName)
}

// exportCustomQuery executa query livre e exporta
func exportCustomQuery() {
	fmt.Println("\n📝 Digite sua query SQL (termine com ;):")
	fmt.Print("> ")

	var queryLines []string
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		queryLines = append(queryLines, line)
		if strings.HasSuffix(line, ";") {
			break
		}
		fmt.Print("  ")
	}

	query := strings.Join(queryLines, " ")
	query = strings.TrimSuffix(query, ";")

	fmt.Println("🔍 Executando query...")
	headers, rows, err := database.QueryCustom(query)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	fmt.Printf("📊 %d colunas | %d registros encontrados\n", len(headers), len(rows))
	doExport(headers, rows, "query_resultado")
}

// doExport pergunta o formato e exporta
func doExport(headers []string, rows [][]string, baseName string) {
	if len(rows) == 0 {
		fmt.Println("⚠️  Nenhum dado para exportar.")
		return
	}

	fmt.Println("\n📤 Formato de exportação:")
	fmt.Println("  1. CSV")
	fmt.Println("  2. Excel (.xlsx)")
	fmt.Println("  3. PDF")
	fmt.Println("  4. Todos os formatos")

	op := input("Escolha: ")

	timestamp := time.Now().Format("20060102_150405")
	base := fmt.Sprintf("%s_%s", baseName, timestamp)

	formats := map[string]string{}
	switch op {
	case "1":
		formats["csv"] = base + ".csv"
	case "2":
		formats["excel"] = base + ".xlsx"
	case "3":
		formats["pdf"] = base + ".pdf"
	case "4":
		formats["csv"] = base + ".csv"
		formats["excel"] = base + ".xlsx"
		formats["pdf"] = base + ".pdf"
	default:
		fmt.Println("❌ Opção inválida.")
		return
	}

	fmt.Println()
	for format, filename := range formats {
		if err := exporter.Export(format, headers, rows, filename); err != nil {
			fmt.Printf("❌ Erro ao exportar %s: %v\n", format, err)
		}
	}
}

func inputDefault(prompt, def string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return def
	}
	return text
}

func printBanner() {
	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════╗")
	fmt.Println("  ║       AUTOMAÇÃO GO - EXPORTADOR      ║")
	fmt.Println("  ║   MySQL → CSV | Excel | PDF          ║")
	fmt.Println("  ╚══════════════════════════════════════╝")
	fmt.Println()
}
