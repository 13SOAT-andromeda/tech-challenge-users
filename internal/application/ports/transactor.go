package ports

type Repositories struct {
	Person   PersonRepository
	User     UserRepository
	Employee EmployeeRepository
	Customer CustomerRepository
}

type Transactor interface {
	WithinTransaction(fn func(repos Repositories) error) error
}
