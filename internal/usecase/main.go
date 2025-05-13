package usecase

type UseCase struct {
	MovieRepo MovieRepoI
	ActorRepo ActorRepoI
}

func NewUseCase(
	movieRepo MovieRepoI,
	actorRepo ActorRepoI,

) *UseCase {
	return &UseCase{
		MovieRepo: movieRepo,
		ActorRepo: actorRepo,
	}
}
