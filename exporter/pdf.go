package exporter

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

type PDFExporter struct{}

func (e *PDFExporter) Export(headers []string, rows [][]string, filename string) error {
	if filename == "" {
		filename = "export.pdf"
	}

	pdf := gofpdf.New("L", "mm", "A4", "") // Landscape para tabelas largas
	pdf.SetMargins(10, 15, 10)
	pdf.AddPage()

	// Tradutor para suportar acentos (UTF-8 -> Latin-1)
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// Título
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(45, 106, 159)
	pdf.CellFormat(0, 10, tr("Relatório de Dados"), "", 1, "C", false, 0, "")
	pdf.Ln(3)

	// Calcula largura das colunas
	pageWidth, _ := pdf.GetPageSize()
	margins, _, margins2, _ := pdf.GetMargins()
	usableWidth := pageWidth - margins - margins2
	colWidth := usableWidth / float64(len(headers))
	if colWidth > 60 {
		colWidth = 60
	}
	rowHeight := 7.0

	// Cabeçalho da tabela
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(45, 106, 159)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetDrawColor(200, 200, 200)

	for _, h := range headers {
		label := h
		if len(label) > 18 {
			label = label[:16] + ".."
		}
		pdf.CellFormat(colWidth, rowHeight, tr(label), "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Linhas de dados
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(50, 50, 50)

	for i, row := range rows {
		// Zebrado
		if i%2 == 0 {
			pdf.SetFillColor(235, 243, 251)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		// Nova página se necessário
		if pdf.GetY() > 185 {
			pdf.AddPage()
			// Repete cabeçalho
			pdf.SetFont("Arial", "B", 9)
			pdf.SetFillColor(45, 106, 159)
			pdf.SetTextColor(255, 255, 255)
			for _, h := range headers {
				label := h
				if len(label) > 18 {
					label = label[:16] + ".."
				}
				pdf.CellFormat(colWidth, rowHeight, tr(label), "1", 0, "C", true, 0, "")
			}
			pdf.Ln(-1)
			pdf.SetFont("Arial", "", 8)
			pdf.SetTextColor(50, 50, 50)
		}

		for j, val := range row {
			if j >= len(headers) {
				break
			}
			label := val
			if len(label) > 18 {
				label = label[:16] + ".."
			}
			pdf.CellFormat(colWidth, rowHeight, tr(label), "1", 0, "L", true, 0, "")
		}
		pdf.Ln(-1)
	}

	// Rodapé com total de registros
	pdf.Ln(3)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, 6, tr(fmt.Sprintf("Total de registros: %d", len(rows))), "", 0, "R", false, 0, "")

	if err := pdf.OutputFileAndClose(filename); err != nil {
		return fmt.Errorf("erro ao salvar PDF: %w", err)
	}

	fmt.Printf("✅ PDF exportado: %s (%d linhas)\n", filename, len(rows))
	return nil
}
