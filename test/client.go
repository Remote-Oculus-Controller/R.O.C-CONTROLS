package roc_test

type client interface {
	Start()
	Stop()
	Send()
}
