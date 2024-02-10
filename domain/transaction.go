package domain

type Transaction struct {
	id           int
	client_id    int
	valor        int
	descricao    string
	realizada_em string
}
