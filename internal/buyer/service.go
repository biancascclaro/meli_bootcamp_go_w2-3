package buyer

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-3/internal/domain"
)

// Errors
var (
	ErrNotFound = errors.New("buyer not found")
	ErrExists = errors.New("buyer already exists")
)

type Service interface{	
	GetAll(ctx context.Context) ([]domain.Buyer, error)
	Get(ctx context.Context, id int) (domain.Buyer, error)
	Save(ctx context.Context, b domain.Buyer) (int, error)
}

type buyerService struct{
	repository Repository
}

func NewService(r Repository) Service {
	return &buyerService{
		repository : r,
	}
}

func (b *buyerService) GetAll(ctx context.Context) ([]domain.Buyer, error) {
	buyers, err := b.repository.GetAll(ctx)
	if err != nil {
		return buyers, err
	}
	return buyers, nil
}

func (b *buyerService) Get(ctx context.Context, id int) (domain.Buyer, error) {
	buyer, err := b.repository.Get(ctx, id)
	if err != nil {
		if err.Error()=="sql: no rows in result set"{
			return domain.Buyer{}, ErrNotFound
		}
		return domain.Buyer{}, err
	}
	return buyer, err
}

func (b *buyerService) Save(ctx context.Context, d domain.Buyer) (int, error){
	userExist := b.repository.Exists(ctx, d.CardNumberID)
	if userExist{
		return 0, ErrExists
	}
	sellerId, err := b.repository.Save(ctx, d)
	return sellerId, err
}