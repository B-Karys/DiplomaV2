package usecase

type TokenUseCase interface {
	CreateActivationToken(email string) error
	CreateAuthenticationToken(email, password string) (string, error)
	CreatePasswordResetToken(email string) error
}
