package usecase

type Usecase struct {
	dbrepo    DBRepo
	redisrepo RedisRepo
}

// NewURLUsecase creates a new instance of URLUsecase with the given repository.
func NewUsecase(dbrepo DBRepo, redisrepo RedisRepo) *Usecase {
	return &Usecase{
		dbrepo:    dbrepo,
		redisrepo: redisrepo,
	}
}
