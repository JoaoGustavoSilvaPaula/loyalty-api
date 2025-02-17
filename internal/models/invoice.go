package models

type Invoice struct {
	InvoiceNumber string `json:"invoice_number"`
	IssueDate     string `json:"issue_date"`
	CNPJ          string `json:"cnpj"`
}
