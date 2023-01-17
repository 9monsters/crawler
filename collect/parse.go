package collect

type Rule struct {
	ParseFunc func(*Context) ParseResult
}

type RuleTree struct {
	Root  func() []*Request
	Trunk map[string]*Rule
}
