package converters

import (
	usermodel "tech-challenge-users/internal/adapter/database/model/user"
	"tech-challenge-users/internal/domain"
)

func UserToModel(u domain.User) usermodel.Model {
	return usermodel.Model{
		ID:       u.ID,
		Password: u.Password,
		Role:     u.Role,
		PersonID: u.PersonID,
	}
}

func UserToDomain(m usermodel.Model) domain.User {
	return domain.User{
		ID:        m.ID,
		Password:  m.Password,
		Role:      m.Role,
		PersonID:  m.PersonID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
