package context

type AppContext struct {
	ErrorChannel chan<- error
	DoneChannel  <-chan struct{}
}
