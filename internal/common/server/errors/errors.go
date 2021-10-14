package errors

type Common struct {
	Msg string	`json:"msg"`
}

func (c Common) Error() string {
	return c.Msg
}
