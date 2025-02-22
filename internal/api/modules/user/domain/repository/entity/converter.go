package entity

import "github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain"

func ConvertDomainToEntity(domain domain.UserDomainInterface) *User {
	return &User{
		FirstName:  domain.GetFirstName(),
		LastName:   domain.GetLastName(),
		Email:      domain.GetEmail(),
		Password:   domain.GetPassword(),
		IsVerified: domain.GetIsVerified(),
		CreatedAt:  domain.GetCreatedAt(),
		UpdatedAt:  domain.GetUpdatedAt(),
	}
}

func ConvertEntityToDomain(entity *User) domain.UserDomainInterface {
	domain := domain.NewUserDomain(
		entity.FirstName,
		entity.LastName,
		entity.Email,
		entity.Password,
	)

	domain.SetID(entity.ID.Hex())
	domain.SetCreatedAt(entity.CreatedAt)
	domain.SetUpdatedAt(entity.UpdatedAt)
	domain.SetIsVerified(entity.IsVerified)

	return domain
}
