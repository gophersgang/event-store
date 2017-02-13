package vstore

import (
	"errors"
)

// Transaction is any type capable of saving a model to VStore
type Transaction interface {
	Save(Model) error
}

type transaction struct {
	toSave Model
}

var (
	//ErrInvalidTx is raised when we detect an invalid transaction - when nothing is being saved or when transactions are being nested
	ErrInvalidTx = errors.New("Invalid transaction.")
)

//Save sets the provided model to be saved to VStore
func (t *transaction) Save(m Model) error {
	if (t.toSave != nil) {
		return ErrInvalidTx
	}
	t.toSave = m
	return nil
}
