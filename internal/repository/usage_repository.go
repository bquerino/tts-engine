package repository

type UsageRepository interface {
	IncrementMessageCount(provider string) error
	IncrementCharacterCount(provider string, count int) error
	GetUsage(provider string) (messages int, characters int, err error)
	ResetUsage(provider string) error
}
