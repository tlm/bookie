package controller

type Controller interface {
	ActionStream() chan<- Action
}
