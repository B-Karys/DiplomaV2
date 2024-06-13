package usecase

import (
	"DiplomaV2/backend/internal/entity"
	"DiplomaV2/backend/post"
	"DiplomaV2/backend/post/repository"
	"github.com/pkg/errors"
)

type postUseCaseImpl struct {
	Repo repository.PostRepository
}

var (
	ErrorFailedPostValidation = errors.New("Post doesn't belong to this user")
)

func NewPostUseCase(repository repository.PostRepository) PostUseCase {
	return &postUseCaseImpl{
		Repo: repository,
	}
}

func (p *postUseCaseImpl) GetPostById(id int64) (*entity.Post, error) {
	thePost, err := p.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return thePost, nil
}

func (p *postUseCaseImpl) CreatePost(post *entity.Post) error {
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

func (p *postUseCaseImpl) UpdatePost(postID, userID int64, updatedPost *entity.Post) error {
	thePost, err := p.Repo.GetByID(postID)
	if err != nil {
		return err
	}

	if thePost.AuthorID != userID {
		return ErrorFailedPostValidation
	}

	thePost.Name = updatedPost.Name
	thePost.Description = updatedPost.Description
	thePost.Type = updatedPost.Type
	thePost.Skills = updatedPost.Skills
	thePost.Version += 1

	err = p.Repo.Update(thePost)
	if err != nil {
		return err
	}

	return nil
}

func (p *postUseCaseImpl) GetFilteredPosts(post *entity.Post, filters postsFilter.Filters) ([]*entity.Post, postsFilter.Metadata, error) {
	filteredPosts, metadata, err := p.Repo.GetFilteredPosts(post, filters)
	if err != nil {
		return nil, metadata, err
	}
	return filteredPosts, metadata, nil
}
