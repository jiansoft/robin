package robin

type Disposable interface {
	Dispose()
	Identify() string
}
