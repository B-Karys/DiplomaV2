package usecase

type TokenUseCase interface {
	createActivationToken(email string) error
	createAuthenticationToken(email, password string) (string, error)
	createPasswordResetToken(email string) error
}
