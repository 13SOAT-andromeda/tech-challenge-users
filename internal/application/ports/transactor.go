package ports

type Repositories struct {
	Person   PersonRepository
	User     UserRepository
	Employee EmployeeRepository
}

type Transactor interface {
	WithinTransaction(fn func(repos Repositories) error) error
}
