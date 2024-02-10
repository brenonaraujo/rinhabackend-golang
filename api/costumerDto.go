package api

type CostumerDto struct {
	Valor     int    `json:"valor" binding:"required"`
	Tipo      string `json:"tipo" binding:"required"`
	Descricao string `json:"descricao" binding:"required"`
}
