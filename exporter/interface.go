package exporter

import "fmt"

// Exporter é a interface que todos os exportadores devem implementar
type Exporter interface {
	Export(headers []string, rows [][]string, filename string) error
}

// Export escolhe o exportador correto com base no formato
func Export(format string, headers []string, rows [][]string, filename string) error {
	var exp Exporter

	switch format {
	case "csv":
		exp = &CSVExporter{}
	case "excel":
		exp = &ExcelExporter{}
	case "pdf":
		exp = &PDFExporter{}
	default:
		return fmt.Errorf("formato não suportado: %s (use: csv, excel, pdf)", format)
	}

	return exp.Export(headers, rows, filename)
}
