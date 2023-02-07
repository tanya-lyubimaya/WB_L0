package domain

type OrderRepository interface {
	Create(Order) (*Order, error)
	Read(Order, int) ([]*Order, error)
	ReadByID(Order) (*Order, error)
	Update(Order) (*Order, error)
	Delete(Order) error
	Close()
}

type UseCase interface {
	ReadOrderByID(string) (*Order, error)
	Close()
}
