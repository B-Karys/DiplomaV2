package usecase

import (
	"DiplomaV2/domain/models"
	"DiplomaV2/repositories"
)

type postUseCaseImpl struct {
	Repo repositories.PostRepository
}

func NewPostUseCase(repository repositories.PostRepository) PostUseCase {
	return &postUseCaseImpl{
		Repo: repository,
	}
}

func (pos *postUseCaseImpl) CreatePost(post *models.Post) error {
	err := pos.Repo.Insert(post)
	if err != nil {
		return err
	}
	return nil
}

func (pos *postUseCaseImpl) ShowOncePost(id int64) (*models.Post, error) {
	post, err := pos.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (pos *postUseCaseImpl) DeletePost(id int64) error {
	err := pos.Repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (pos *postUseCaseImpl) UpdatePost(post *models.Post) error {
	err := pos.Repo.Update(post)
	if err != nil {
		return err
	}
	return nil
}
