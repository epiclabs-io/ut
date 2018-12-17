package ut

type serviceWrapper struct {
	close func()
}

func NewService(cleanup func()) *serviceWrapper {
	return &serviceWrapper{
		close: cleanup,
	}
}

func (sw *serviceWrapper) Close() {
	sw.close()
}
