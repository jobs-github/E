package function

type Function struct {
	String     func() string
	ArgumentOf func(idx int) string
	Body       func() string
}
