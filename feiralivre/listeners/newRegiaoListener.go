package listeners

type RegiaoListener struct {
}

func (fl RegiaoListener) Topic() string {
	return "regiao"
}

func (fl RegiaoListener) Subscription() string {
	return "newRegiaoSubcriber"
}

func (fl RegiaoListener) URL() string {
	return "consumers/novaregiao"
}

func (fl RegiaoListener) IsShared() bool {
	return true
}
