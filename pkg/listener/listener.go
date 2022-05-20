package listener

type Listener interface {
	RunListenerService()
	StopService()
}
