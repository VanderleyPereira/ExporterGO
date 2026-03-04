package exporter

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ExcelExporter struct{}

func (e *ExcelExporter) Export(headers []string, rows [][]string, filename string) error {
	if filename == "" {
		filename = "export.xlsx"
	}

	f := excelize.NewFile()
	sheet := "Dados"
	f.SetSheetName("Sheet1", sheet)

	// Estilo do cabeçalho: fundo azul escuro, texto branco, negrito
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "#FFFFFF",
			Size:  11,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#2D6A9F"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#FFFFFF", Style: 1},
			{Type: "right", Color: "#FFFFFF", Style: 1},
			{Type: "bottom", Color: "#FFFFFF", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("erro ao criar estilo: %w", err)
	}

	// Estilo das linhas pares (zebrado)
	evenStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#EBF3FB"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})

	// Escreve cabeçalhos
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// Escreve linhas
	for rowIdx, row := range rows {
		for colIdx, val := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, val)
			// Zebrado nas linhas pares
			if rowIdx%2 == 1 {
				f.SetCellStyle(sheet, cell, cell, evenStyle)
			}
		}
	}

	// Auto-ajusta largura das colunas
	for i, h := range headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		maxLen := len(h)
		for _, row := range rows {
			if i < len(row) && len(row[i]) > maxLen {
				maxLen = len(row[i])
			}
		}
		width := float64(maxLen) + 4
		if width > 50 {
			width = 50
		}
		f.SetColWidth(sheet, col, col, width)
	}

	// Congela a primeira linha (cabeçalho)
	f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	if err := f.SaveAs(filename); err != nil {
		return fmt.Errorf("erro ao salvar Excel: %w", err)
	}

	fmt.Printf("✅ Excel exportado: %s (%d linhas)\n", filename, len(rows))
	return nil
}
