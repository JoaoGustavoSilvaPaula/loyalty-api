package services

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/joaogustavosp/loyalty-api/internal/models"
)

type ExtractService struct{}

func NewExtractService() *ExtractService {
	return &ExtractService{}
}

func (s *ExtractService) ExtractData(url string) ([]models.ProductInvoice, models.Invoice, error) {
	// Realizar a requisição HTTP
	resp, err := http.Get(url)
	if err != nil {
		return nil, models.Invoice{}, err
	}
	defer resp.Body.Close()

	// Verificar se a requisição foi bem-sucedida
	if resp.StatusCode != 200 {
		return nil, models.Invoice{}, fmt.Errorf("Erro ao acessar URL: %s", resp.Status)
	}

	// Parsear o HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, models.Invoice{}, err
	}

	var productInvoices []models.ProductInvoice
	var invoice models.Invoice

	// Expressão regular para extrair apenas números, pontos e vírgulas
	reNumber := regexp.MustCompile(`[0-9,.]+`)
	// Expressão regular para extrair o código
	reCode := regexp.MustCompile(`Código: (\d+)`)
	// Expressão regular para extrair o CNPJ
	reCNPJ := regexp.MustCompile(`CNPJ: (\d{2}\.\d{3}\.\d{3}\/\d{4}\-\d{2})`)

	// Encontrar o elemento com id="collapse4" e extrair os dados da nota
	collapse4 := doc.Find("#collapse4")
	invoice.InvoiceNumber = strings.TrimSpace(collapse4.Find("table:nth-child(8) > tbody > tr > td:nth-child(3)").Text())
	invoice.IssueDate = strings.TrimSpace(collapse4.Find("table:nth-child(8) > tbody > tr > td:nth-child(4)").Text())
	cnpjText := doc.Find("td:contains('CNPJ')").Text()
	cnpj := reCNPJ.FindStringSubmatch(cnpjText)
	if len(cnpj) > 1 {
		invoice.CNPJ = cnpj[1]
	} else {
		invoice.CNPJ = ""
	}

	// Selecionar e extrair os dados dos produtos
	doc.Find("table.table-striped tbody tr").Each(func(i int, s *goquery.Selection) {
		codeText := s.Find("td:nth-child(1)").Text()
		code := reCode.FindStringSubmatch(codeText)
		if len(code) > 1 {
			codeText = code[1]
		} else {
			codeText = ""
		}

		ProductInvoice := models.ProductInvoice{
			Name:     strings.TrimSpace(s.Find("td:nth-child(1) h7").Text()),
			Code:     codeText,
			Quantity: reNumber.FindString(s.Find("td:nth-child(2)").Text()),
			Unit:     strings.TrimSpace(s.Find("td:nth-child(3)").Text()),
			Value:    reNumber.FindString(s.Find("td:nth-child(4)").Text()),
		}
		productInvoices = append(productInvoices, ProductInvoice)
	})

	return productInvoices, invoice, nil
}
