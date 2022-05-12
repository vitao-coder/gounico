package listeners

type FeiraListener struct {
}

func (fl FeiraListener) Topic() string {
	return "feira"
}

func (fl FeiraListener) Subscription() string {
	return "newfeiraSubcriber"
}
func (fl FeiraListener) URL() string {
	return "consumers/novafeira"
}

func (fl FeiraListener) IsShared() bool {
	return false
}
