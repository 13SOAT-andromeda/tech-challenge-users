package services

import (
	"errors"

	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
	"tech-challenge-users/pkg/encryption"
)

var ErrCustomerNotFound = errors.New("customer not found")

type CreateCustomerInput struct {
	Name     string
	Email    string
	Password string
	Type     string
	Document string
	Contact  string
	Address  domain.Address
}

type UpdateCustomerInput struct {
	Name    *string
	Email   *string
	Contact *string
	Address *domain.Address
	Type    *string
}

type CustomerOutput struct {
	Customer domain.Customer
	Person   domain.Person
}

type CustomerService struct {
	transactor   ports.Transactor
	personRepo   ports.PersonRepository
	customerRepo ports.CustomerRepository
}

func NewCustomerService(
	transactor ports.Transactor,
	personRepo ports.PersonRepository,
	customerRepo ports.CustomerRepository,
) *CustomerService {
	return &CustomerService{
		transactor:   transactor,
		personRepo:   personRepo,
		customerRepo: customerRepo,
	}
}

func (s *CustomerService) createCustomerInTx(repos ports.Repositories, input CreateCustomerInput) (CustomerOutput, error) {
	existing, err := repos.Person.GetByEmail(input.Email)
	if err != nil {
		return CustomerOutput{}, err
	}
	if existing != nil {
		return CustomerOutput{}, ErrEmailAlreadyExists
	}

	existingDoc, err := repos.Person.GetByDocument(input.Document)
	if err != nil {
		return CustomerOutput{}, err
	}
	if existingDoc != nil {
		return CustomerOutput{}, ErrDocumentAlreadyExists
	}

	person := domain.Person{
		Name:     input.Name,
		Email:    input.Email,
		Contact:  input.Contact,
		Document: input.Document,
		IsActive: true,
		Address:  input.Address,
	}
	if err := repos.Person.Create(&person); err != nil {
		return CustomerOutput{}, err
	}

	hash, err := encryption.Hash(input.Password)
	if err != nil {
		return CustomerOutput{}, err
	}

	user := domain.User{
		Password: hash,
		Role:     domain.RoleCustomer,
		PersonID: person.ID,
	}
	if err := repos.User.Create(&user); err != nil {
		return CustomerOutput{}, err
	}

	customer := domain.Customer{
		Type:     input.Type,
		PersonID: person.ID,
	}
	if err := repos.Customer.Create(&customer); err != nil {
		return CustomerOutput{}, err
	}

	return CustomerOutput{Customer: customer, Person: person}, nil
}

func (s *CustomerService) CreateCustomer(input CreateCustomerInput) (*CustomerOutput, error) {
	if err := domain.ValidateCustomerType(input.Type); err != nil {
		return nil, err
	}
	if _, err := domain.NewPassword(input.Password); err != nil {
		return nil, err
	}
	if _, err := domain.NewDocument(input.Document); err != nil {
		return nil, err
	}

	var output CustomerOutput
	err := s.transactor.WithinTransaction(func(repos ports.Repositories) error {
		var err error
		output, err = s.createCustomerInTx(repos, input)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &output, nil
}

func (s *CustomerService) GetByID(id int64) (*CustomerOutput, error) {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, ErrCustomerNotFound
	}

	person, err := s.personRepo.FindByID(customer.PersonID)
	if err != nil {
		return nil, err
	}

	return &CustomerOutput{Customer: *customer, Person: *person}, nil
}

func (s *CustomerService) ListCustomers(filters ports.CustomerFilters) ([]CustomerOutput, error) {
	customers, err := s.customerRepo.FindAll(filters)
	if err != nil {
		return nil, err
	}

	outputs := make([]CustomerOutput, 0, len(customers))
	for _, c := range customers {
		person, err := s.personRepo.FindByID(c.PersonID)
		if err != nil {
			return nil, err
		}
		if person == nil {
			continue
		}
		outputs = append(outputs, CustomerOutput{Customer: c, Person: *person})
	}
	return outputs, nil
}

func (s *CustomerService) applyEmailUpdate(person *domain.Person, email string) error {
	existing, err := s.personRepo.GetByEmail(email)
	if err != nil {
		return err
	}
	if existing != nil && existing.ID != person.ID {
		return ErrEmailAlreadyExists
	}
	person.Email = email
	return nil
}

func (s *CustomerService) applyTypeUpdate(customer *domain.Customer, t string) error {
	if err := domain.ValidateCustomerType(t); err != nil {
		return err
	}
	customer.Type = t
	return s.customerRepo.Update(customer)
}

func (s *CustomerService) UpdateCustomer(id int64, input UpdateCustomerInput) (*CustomerOutput, error) {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, ErrCustomerNotFound
	}

	person, err := s.personRepo.FindByID(customer.PersonID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		person.Name = *input.Name
	}
	if input.Email != nil {
		if err := s.applyEmailUpdate(person, *input.Email); err != nil {
			return nil, err
		}
	}
	if input.Contact != nil {
		person.Contact = *input.Contact
	}
	if input.Address != nil {
		person.Address = *input.Address
	}
	if err := s.personRepo.Update(person); err != nil {
		return nil, err
	}

	if input.Type != nil {
		if err := s.applyTypeUpdate(customer, *input.Type); err != nil {
			return nil, err
		}
	}

	return &CustomerOutput{Customer: *customer, Person: *person}, nil
}

func (s *CustomerService) DeleteCustomer(id int64) error {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return err
	}
	if customer == nil {
		return ErrCustomerNotFound
	}
	return s.customerRepo.Delete(id)
}
