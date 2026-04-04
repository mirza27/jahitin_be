package service

import (
	"context"
	"database/sql"
	"errors"

	db "jahitin_be/database/repository"
)

type CustomerService interface {
	ListUserCustomers(ctx context.Context, filter ListCustomersFilter) ([]*db.Customer, error)
	GetUserCustomerDetail(ctx context.Context, user_id int64, customer_id int64) (*db.Customer, error)
}

type customerService struct {
	store db.Store
}

func NewCustomerService(store db.Store) CustomerService {
	return &customerService{store: store}
}

type ListCustomersFilter struct {
	UserID int64
	Search string
	Page   int32
	Limit  int32
}

func (h *customerService) ListUserCustomers(ctx context.Context, filter ListCustomersFilter) ([]*db.Customer, error) {

	var customers []db.Customer
	var err error

	if filter.Search == "" {
		customers, err = h.store.ListCustomersByUserID(ctx, db.ListCustomersByUserIDParams{
			UserID: filter.UserID,
			Limit:  filter.Limit,
			Offset: (filter.Page - 1) * filter.Limit,
		})
	} else {
		customers, err = h.store.ListCustomersByUserIDAndName(ctx, db.ListCustomersByUserIDAndNameParams{
			UserID:  filter.UserID,
			Column2: sql.NullString{String: filter.Search, Valid: true},
			Limit:   filter.Limit,
			Offset:  (filter.Page - 1) * filter.Limit,
		})
	}
	if err != nil {
		return nil, err
	}

	customerPtrs := make([]*db.Customer, len(customers))
	for i := range customers {
		customerPtrs[i] = &customers[i]
	}

	return customerPtrs, err
}

func (h *customerService) GetUserCustomerDetail(ctx context.Context, user_id int64, customer_id int64) (*db.Customer, error) {

	customer, err := h.store.GetCustomerByID(ctx, customer_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("customer not found")
		}

		return nil, err
	}

	if customer.UserID != user_id {
		return nil, errors.New("unauthorized access to customer data")
	}

	return &customer, nil

}
