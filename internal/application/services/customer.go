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
		existing, err := repos.Person.GetByEmail(input.Email)
		if err != nil {
			return err
		}
		if existing != nil {
			return ErrEmailAlreadyExists
		}

		existingDoc, err := repos.Person.GetByDocument(input.Document)
		if err != nil {
			return err
		}
		if existingDoc != nil {
			return ErrDocumentAlreadyExists
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
			return err
		}

		hash, err := encryption.Hash(input.Password)
		if err != nil {
			return err
		}

		user := domain.User{
			Password: hash,
			Role:     domain.RoleCustomer,
			PersonID: person.ID,
		}
		if err := repos.User.Create(&user); err != nil {
			return err
		}

		customer := domain.Customer{
			Type:     input.Type,
			PersonID: person.ID,
		}
		if err := repos.Customer.Create(&customer); err != nil {
			return err
		}

		output = CustomerOutput{Customer: customer, Person: person}
		return nil
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
		existing, err := s.personRepo.GetByEmail(*input.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != person.ID {
			return nil, ErrEmailAlreadyExists
		}
		person.Email = *input.Email
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
		if err := domain.ValidateCustomerType(*input.Type); err != nil {
			return nil, err
		}
		customer.Type = *input.Type
		if err := s.customerRepo.Update(customer); err != nil {
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
