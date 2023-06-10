package stderr

import help "github.com/mdanialr/pwman_backend/pkg/helper"

// NewUC return pointer to UC using given code and message as the constructor.
func NewUC(c, m string) error {
	return &UC{Code: c, Msg: m}
}

// UC standard error object that may be returned by use case layer.
type UC struct {
	Code string
	Msg  string
}

// Error implement error interface.
func (u *UC) Error() string {
	return help.Pad(u.Code+":", u.Msg)
}
