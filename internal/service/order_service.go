package service

import (
	"context"
	"database/sql"
	"errors"
	db "jahitin_be/database/repository"
	"time"
)

type OrderService interface {
	CreateUserOrder(ctx context.Context, user_id int64, o_data CreateUserOrderInput, c_data CustomerInput) (bool, error)
}

type orderService struct {
	store db.Store
}

func NewOrderService(store db.Store) OrderService {
	return &orderService{store: store}
}

type OrderItemInput struct {
	ClothesFor          string
	Notes               string
	ClothesCategoryID   int64
	ServiceTypeID       int64
	CustomServiceName   string
	Price               int64
	IsSaveCustomerNotes bool
}
type CreateUserOrderInput struct {
	Name     string
	Deadline *time.Time
	Items    []OrderItemInput
}

type CustomerInput struct {
	IsNewCustomer bool
	CustomerID    int64
	Name          string
	Phone         string
}

func (o *orderService) CreateUserOrder(ctx context.Context, user_id int64, o_data CreateUserOrderInput, c_data CustomerInput) (bool, error) {
	err := o.store.ExecTx(ctx, func(q *db.Queries) error {

		var customerID int64

		// if new customer
		if c_data.IsNewCustomer {
			cParam := db.CreateCustomerParams{
				Name:   c_data.Name,
				UserID: user_id,
				Phone:  sql.NullString{String: c_data.Phone, Valid: c_data.Phone != ""},
				Notes:  sql.NullString{String: "", Valid: false},
			}

			customer, err := q.CreateCustomer(ctx, cParam)
			if err != nil {
				return err
			}

			customerID = customer.ID // replace
		}

		// check if this user's customer
		if !c_data.IsNewCustomer {
			customer, err := o.store.GetCustomerByID(ctx, c_data.CustomerID)
			if err != nil {
				if err == sql.ErrNoRows {
					return errors.New("customer not found")
				}

			}

			if customer.UserID != user_id {
				return errors.New("customer does not belong to user")
			}
		}

		oParam := db.CreateOrderParams{
			UserID:     user_id,
			CustomerID: customerID,
			Name:       o_data.Name,
			Deadline:   sql.NullTime{Time: timeOrZero(o_data.Deadline), Valid: o_data.Deadline != nil},
		}

		order, err := q.CreateOrder(ctx, oParam)
		if err != nil {
			return err
		}

		// create each order item
		for _, item := range o_data.Items {
			serviceName := sql.NullString{Valid: false}
			serviceTypeID := sql.NullInt64{Int64: item.ServiceTypeID, Valid: item.ServiceTypeID != 0}
			if item.CustomServiceName != "" {
				serviceName = sql.NullString{String: item.CustomServiceName, Valid: true}
				serviceTypeID = sql.NullInt64{Valid: false}
			}

			oiParam := db.CreateOrderItemParams{
				OrderID:           order.ID,
				CategoryID:        sql.NullInt64{Int64: item.ClothesCategoryID, Valid: item.ClothesCategoryID != 0},
				ServiceTypeID:     serviceTypeID,
				ClothesFor:        item.ClothesFor,
				CustomServiceName: serviceName,
				Notes:             sql.NullString{String: item.Notes, Valid: item.Notes != ""},
				Price:             item.Price,
			}

			if _, err := q.CreateOrderItem(ctx, oiParam); err != nil {
				return err
			}

			// if user wants save note to this customer
			if item.IsSaveCustomerNotes {
				if err := o.saveCustomerOrderNotes(ctx, q, item.Notes, customerID); err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (o *orderService) saveCustomerOrderNotes(ctx context.Context, q db.Querier, notes string, customer_id int64) error {

	_, err := q.UpdateCustomerNotes(ctx, db.UpdateCustomerNotesParams{
		ID:    customer_id,
		Notes: sql.NullString{String: notes, Valid: notes != ""},
	})
	if err != nil {
		return err
	}

	return nil
}

func timeOrZero(value *time.Time) time.Time {
	if value == nil {
		return time.Time{}
	}

	return *value
}
