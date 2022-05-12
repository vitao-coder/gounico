package domain

type GenericListener struct {
	topic        string
	subscription string
	url          string
	isShared     bool
}

func NewGenericListener(topic string, subscription string, url string, isShared bool) GenericListener {
	return GenericListener{
		topic:        topic,
		subscription: subscription,
		url:          url,
		isShared:     isShared,
	}

}

func (gl GenericListener) Topic() string {
	return gl.topic
}

func (gl GenericListener) Subscription() string {
	return gl.subscription
}

func (gl GenericListener) URL() string {
	return gl.url
}

func (gl GenericListener) IsShared() bool {
	return gl.isShared
}
