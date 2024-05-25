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

func (p *postUseCaseImpl) GetPostById(id int64) (*models.Post, error) {
	post, err := p.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

/*
func (p *postUseCaseImpl) GetPostByAuthor(authorID int64) ([]*models.Post, error) {
	posts, err := p.Repo.GetByAuthor(authorID)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
*/

func (p *postUseCaseImpl) GetAllPosts() ([]*models.Post, error) {
	posts, err := p.Repo.GetAll()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

/*
func (p *postUseCaseImpl) GetPostByType(postType string) ([]*models.Post, error) {
	posts, err := p.Repo.GetByType(postType)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
*/

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
		return err
	}

	post.Name = name
	post.Description = description
	post.Type = postType
	post.Skills = skills

	err = p.Repo.Update(post)
	if err != nil {
		return err
	}

	return err
}

func (p *postUseCaseImpl) GetFilteredPosts(authorID, postType string) ([]*models.Post, error) {
	// Initialize an empty slice to store filtered posts
	var filteredPosts []*models.Post
	var err error

	// If both authorID and postType are provided, filter by both
	if authorID != "" && postType != "" {
		filteredPosts, err = p.Repo.GetByAuthorAndType(authorID, postType)
	} else if authorID != "" {
		// If only authorID is provided, filter by author
		filteredPosts, err = p.Repo.GetByAuthorAndType(authorID, "") // Provide an empty postType for filtering by author only
	} else if postType != "" {
		// If only postType is provided, filter by type
		filteredPosts, err = p.Repo.GetByAuthorAndType("", postType) // Provide an empty authorID for filtering by type only
	} else {
		// If no filter criteria provided, return all posts
		filteredPosts, err = p.Repo.GetAll()
	}

	if err != nil {
		return nil, err
	}

	return filteredPosts, nil
}
