package services

import (
	"errors"
	"log"

	"tech-challenge-users/internal/application/ports"
	"tech-challenge-users/internal/domain"
	"tech-challenge-users/pkg/encryption"
)

var (
	ErrEmailAlreadyExists    = errors.New("email already in use")
	ErrDocumentAlreadyExists = errors.New("document already in use")
	ErrUserNotFound          = errors.New("user not found")
)

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
	Role     string
	Document string
	Contact  string
	Position string
	Address  domain.Address
}

type UpdateUserInput struct {
	Name     *string
	Email    *string
	Contact  *string
	Position *string
	Address  *domain.Address
}

type AdminConfig struct {
	Email    string
	Password string
	Document string
}

type UserOutput struct {
	User     domain.User
	Person   domain.Person
	Employee *domain.Employee
}

type UserService struct {
	transactor   ports.Transactor
	personRepo   ports.PersonRepository
	userRepo     ports.UserRepository
	employeeRepo ports.EmployeeRepository
}

func NewUserService(
	transactor ports.Transactor,
	personRepo ports.PersonRepository,
	userRepo ports.UserRepository,
	employeeRepo ports.EmployeeRepository,
) *UserService {
	return &UserService{
		transactor:   transactor,
		personRepo:   personRepo,
		userRepo:     userRepo,
		employeeRepo: employeeRepo,
	}
}

func (s *UserService) CreateUser(input CreateUserInput) (*UserOutput, error) {
	if _, err := domain.NewPassword(input.Password); err != nil {
		return nil, err
	}
	if _, err := domain.NewDocument(input.Document); err != nil {
		return nil, err
	}
	if err := domain.ValidateRole(input.Role); err != nil {
		return nil, err
	}

	var output UserOutput
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
			Role:     input.Role,
			PersonID: person.ID,
		}
		if err := repos.User.Create(&user); err != nil {
			return err
		}

		var employee *domain.Employee
		if input.Role != domain.RoleCustomer {
			emp := domain.Employee{
				Position: input.Position,
				PersonID: person.ID,
			}
			if err := repos.Employee.Create(&emp); err != nil {
				return err
			}
			employee = &emp
		}

		output = UserOutput{User: user, Person: person, Employee: employee}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &output, nil
}

func (s *UserService) GetByID(id int64) (*UserOutput, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	person, err := s.personRepo.FindByID(user.PersonID)
	if err != nil {
		return nil, err
	}

	employee, err := s.employeeRepo.FindByPersonID(user.PersonID)
	if err != nil {
		return nil, err
	}

	return &UserOutput{User: *user, Person: *person, Employee: employee}, nil
}

func (s *UserService) ListUsers(filters ports.UserFilters) ([]UserOutput, error) {
	users, err := s.userRepo.FindAll(filters)
	if err != nil {
		return nil, err
	}

	outputs := make([]UserOutput, 0, len(users))
	for _, u := range users {
		person, err := s.personRepo.FindByID(u.PersonID)
		if err != nil {
			return nil, err
		}
		if person == nil {
			continue
		}
		employee, err := s.employeeRepo.FindByPersonID(u.PersonID)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, UserOutput{User: u, Person: *person, Employee: employee})
	}
	return outputs, nil
}

func (s *UserService) UpdateUser(id int64, input UpdateUserInput) (*UserOutput, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	person, err := s.personRepo.FindByID(user.PersonID)
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

	employee, err := s.employeeRepo.FindByPersonID(user.PersonID)
	if err != nil {
		return nil, err
	}

	return &UserOutput{User: *user, Person: *person, Employee: employee}, nil
}

func (s *UserService) DeleteUser(id int64) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	return s.userRepo.Delete(id)
}

func (s *UserService) GetByDocument(document string) (*UserOutput, error) {
	person, err := s.personRepo.GetByDocument(document)
	if err != nil {
		return nil, err
	}
	if person == nil {
		return nil, ErrUserNotFound
	}

	user, err := s.userRepo.GetByPersonID(person.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &UserOutput{User: *user, Person: *person}, nil
}

func (s *UserService) CreateAdminUser(cfg AdminConfig) {
	existing, err := s.personRepo.GetByEmail(cfg.Email)
	if err != nil {
		log.Printf("CreateAdminUser: error checking existing admin: %v", err)
		return
	}
	if existing != nil {
		return
	}

	if _, err := s.CreateUser(CreateUserInput{
		Name:     "Administrator",
		Email:    cfg.Email,
		Password: cfg.Password,
		Role:     domain.RoleAdministrator,
		Document: cfg.Document,
		Position: "Administrator",
	}); err != nil {
		log.Printf("CreateAdminUser: failed to create admin: %v", err)
	}
}
