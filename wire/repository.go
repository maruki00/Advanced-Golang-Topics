package main

type Repository struct {
	db map[int]any
}

type Service struct {
	repo *Repository
}

func NewRepository(db map[int]any) *Repository {
	return &Repository{
		db: db,
	}
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (svc *Service) GetName() string {
	return "this is the service"
}
