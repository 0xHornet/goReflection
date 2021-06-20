package urlsource

type Literal []string

func (l Literal) GetURLS() (urls []string, err error) {
	return l, nil
}
