package controller

type Controller interface {
	Run(<-chan struct{})
	SetEventHandler(ActionEventHandler)
}
