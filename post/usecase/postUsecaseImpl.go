package usecase

import (
	"DiplomaV2/post/models"
	"DiplomaV2/post/repository"
	"github.com/pkg/errors"
)

type postUseCaseImpl struct {
	Repo repository.PostRepository
}

var (
	ErrorFailedPostFalidation = errors.New("Post doesn't belong to this user")
)

func NewPostUseCase(repository repository.PostRepository) PostUseCase {
	return &postUseCaseImpl{
		Repo: repository,
	}
}

func (p *postUseCaseImpl) GetPostById(id int64) (*models.Post, error) {
	post, err := p.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (p *postUseCaseImpl) CreatePost(post *models.Post) error {
	err := p.Repo.Insert(post)
	if err != nil {
		return err
	}
	return nil
}

func (p *postUseCaseImpl) DeletePost(id int64) error {
	err := p.Repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (p *postUseCaseImpl) UpdatePost(postID, userID int64, name string, description string, skills []string, postType string) error {

	post, err := p.Repo.GetByID(postID)
	if err != nil {
		return err
	}

	if post.AuthorID != userID {
		return ErrorFailedPostFalidation
	}

	post.Name = name
	post.Description = description
	post.Type = postType
	post.Skills = skills
	post.Version = post.Version + 1

	err = p.Repo.Update(post)
	if err != nil {
		return err
	}

	return err
}

func (p *postUseCaseImpl) GetFilteredPosts(name, description, author, postType string, skills []string, filters models.Filters) ([]*models.Post, error) {
	filteredPosts, err := p.Repo.GetFilteredPosts(name, description, author, postType, skills, filters)
	if err != nil {
		return nil, err
	}
	return filteredPosts, nil
}
