package dto

type TestRequest struct {
	Number  int     `json:"number"        validate:"required,min=100,max=599"`
	Message string  `json:"message"       validate:"required"`
	Mes     *string `json:"mes,omitempty" validate:"omitempty"`
}
