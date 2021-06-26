package server


type Handler interface {
	handle(...interface{}) error
}


type SetHandler struct {
	key string
	value string
}