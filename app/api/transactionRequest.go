package api

type TransactionRequest struct {
	Valor     int    `json:"valor" binding:"required,gt=0"`
	Tipo      string `json:"tipo" binding:"required,oneof=c d"`
	Descricao string `json:"descricao" binding:"required,min=1,max=10"`
}
