package services

import (
	"fmt"
	"gfinancer/internal/domain"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type ReportService struct{}

func NewReportService() *ReportService {
	return &ReportService{}
}

func truncate(text string, max int) string {
	if len(text) > max {
		return text[:max-3] + "..."
	}

	return text
}

func (s *ReportService) GeneratePDF(title string, expenses []domain.Report, total string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(190, 10, title, "", 1, "C", false, 0, "")
	pdf.Ln(5)

	for _, exp := range expenses {
		pdf.SetFont("Arial", "", 12)
		pdf.SetTextColor(0, 0, 0)

		cleanDesc := truncate(exp.Description, 30)
		line1 := fmt.Sprintf("%s - %s%s", exp.Date, cleanDesc, exp.Installment)

		pdf.CellFormat(140, 6, line1, "", 0, "L", false, 0, "")
		pdf.CellFormat(50, 6, exp.Value, "", 1, "R", false, 0, "")

		pdf.SetFont("Arial", "I", 10)
		pdf.SetTextColor(128, 128, 128)

		line2 := fmt.Sprintf("Pessoa: %s  |   Cartão: %s", exp.Person, exp.Card)
		pdf.CellFormat(190, 6, line2, "", 1, "L", false, 0, "")

		pdf.Ln(4)
	}

	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(140, 10, "TOTAL DO MES:", "T", 0, "L", false, 0, "")
	pdf.CellFormat(50, 10, total, "T", 1, "R", false, 0, "")

	fileName := fmt.Sprintf("./relatorio_%d.pdf", time.Now().Unix())

	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		return "", fmt.Errorf("Falha ao salvar o PDF: %w", err)
	}

	return fileName, nil
}
