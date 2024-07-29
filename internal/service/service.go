package service

type Service struct {
	Storage storage
}

func New(st storage) *Service {
	s := &Service{
		Storage: st,
	}
	return s
}
