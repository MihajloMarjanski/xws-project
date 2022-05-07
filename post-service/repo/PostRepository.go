package repo

type PostRepository struct {
}

func New() (*PostRepository, error) {
	repo := &PostRepository{}

	return repo, nil
}
