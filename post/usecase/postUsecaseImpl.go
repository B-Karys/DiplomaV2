package usecase

import (
	"DiplomaV2/post/models"
	"DiplomaV2/post/repository"
)

type postUseCaseImpl struct {
	Repo repository.PostRepository
}

func NewPostUseCase(repository repository.PostRepository) PostUseCase {
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

func (pos *postUseCaseImpl) ShowOnePost(id int64) (*models.Post, error) {
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

func (pos *postUseCaseImpl) UpdatePost(postID, userID int64, name string, description string, postType string) error {
	post, err := pos.Repo.GetByID(postID)
	if err != nil {
		return err
	}

	if post.AuthorID != userID {
		return err
	}

	post.Name = name
	post.Description = description
	post.Type = postType

	err = pos.Repo.Update(post)
	if err != nil {
		return err
	}

	return err

}
