package exporter

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CSVExporter struct{}

func (e *CSVExporter) Export(headers []string, rows [][]string, filename string) error {
	if filename == "" {
		filename = "export.csv"
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo CSV: %w", err)
	}
	defer file.Close()

	// BOM para compatibilidade com Excel ao abrir CSV
	file.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escreve cabeçalho
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("erro ao escrever cabeçalho: %w", err)
	}

	// Escreve linhas
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("erro ao escrever linha: %w", err)
		}
	}

	fmt.Printf("✅ CSV exportado: %s (%d linhas)\n", filename, len(rows))
	return nil
}
