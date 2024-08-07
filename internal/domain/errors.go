package domain

import "errors"

var (
	ErrWrongToken  = errors.New("wrong token address")
	ErrWrongAmount = errors.New("wrong amount")
)
