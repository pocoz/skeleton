package root

type BasicService struct{}

type service interface {
	healthy()
}

// healthy ...
func (s *BasicService) healthy() {}
