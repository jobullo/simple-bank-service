package database

type Service[M any] interface {
	Create(model *M) error
	Delete(id int) error
	FetchById(id int) (*M, error)
	List() ([]*M, error)
	Update(mode *M) error
}
