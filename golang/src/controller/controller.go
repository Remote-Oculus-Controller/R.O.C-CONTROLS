package controller

type Controller interface {
	Type() string
	Init()
	Speed() float64
	Direction() float64
}