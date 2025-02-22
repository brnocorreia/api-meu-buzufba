package entity

import "github.com/brnocorreia/api-meu-buzufba/internal/api/modules/user/domain"

func ConvertDomainToEntity(domain domain.UserDomainInterface) *User {
	return &User{
		FirstName:           domain.GetFirstName(),
		LastName:            domain.GetLastName(),
		Email:               domain.GetEmail(),
		Password:            domain.GetPassword(),
		IsVerified:          domain.GetIsVerified(),
		VerificationToken:   domain.GetVerificationToken(),
		VerificationExpires: domain.GetVerificationExpires(),
		EmailVerifiedAt:     domain.GetEmailVerifiedAt(),
		CreatedAt:           domain.GetCreatedAt(),
		UpdatedAt:           domain.GetUpdatedAt(),
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
	domain.SetVerificationToken(entity.VerificationToken)
	domain.SetVerificationExpires(entity.VerificationExpires)
	domain.SetEmailVerifiedAt(entity.EmailVerifiedAt)

	return domain
}
