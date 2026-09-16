package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	db "jahitin_be/database/repository"
	"time"
)

type OrderService interface {
	CreateUserOrder(ctx context.Context, user_id int64, o_data CreateUserOrderInput, c_data CustomerInput) (bool, error)
	ListUserOrders(ctx context.Context, filter ListOrdersFilter) ([]*ListOrdersOutput, error)
	GetOrderDetailsByOrderID(ctx context.Context, userId int64, orderID int64) (*DetailOrderOutput, error)
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

			customerID = customer.ID
		}

		oParam := db.CreateOrderParams{
			UserID:     user_id,
			CustomerID: customerID,
			Name:       o_data.Name,
			Deadline:   sql.NullTime{Time: TimeOrZero(o_data.Deadline), Valid: o_data.Deadline != nil},
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

type ListOrdersFilter struct {
	UserID int64
	Status string
	Search string
	Page   int
	Limit  int
}

type ListOrdersOutput struct {
	OrderID      int64
	Name         string
	Deadline     *time.Time
	CustomerName string
	CustomerId   int64
	Status       string
	Items        []OrderItemOutput
}

type OrderItemOutput struct {
	ClothesFor          string
	ClothesCategoryID   int64
	ClothesCategoryName string
	ServiceTypeID       int64
	ServiceTypeName     string
	CustomServiceName   string
	Price               int64
	Notes               string
	Status              string
}

func (o *orderService) ListUserOrders(ctx context.Context, filter ListOrdersFilter) ([]*ListOrdersOutput, error) {

	var orders []db.ListUserOrdersFilteredRow
	var err error

	// get headers order
	orders, err = o.store.ListUserOrdersFiltered(ctx, db.ListUserOrdersFilteredParams{
		UserID:  filter.UserID,
		Column2: filter.Status,
		Column3: filter.Search,
		Limit:   int32(filter.Limit),
		Offset:  int32((filter.Page - 1) * filter.Limit),
	})
	if err != nil {
		return nil, err
	}

	// get detail order items
	orderIDs := make([]int32, len(orders))
	for i, order := range orders {
		orderIDs[i] = int32(order.ID)
	}

	// get orderItems by orderIDs
	orderItems, err := o.store.ListOrderItemsByMultipleOrderIDs(ctx, orderIDs)
	if err != nil {
		return nil, err
	}

	fmt.Println("orderItems:", orderItems)

	// combine order with order items
	combinedOrders := o.combineOrdersWithItems(orders, orderItems)

	return combinedOrders, nil
}

func (o *orderService) combineOrdersWithItems(orders []db.ListUserOrdersFilteredRow, orderItems []db.ListOrderItemsByMultipleOrderIDsRow) []*ListOrdersOutput {

	// create hashmap orderID -> order items
	orderMap := make(map[int64][]OrderItemOutput)

	for _, item := range orderItems {
		orderMap[item.OrderID] = append(orderMap[item.OrderID], OrderItemOutput{
			ClothesFor:          item.ClothesFor,
			ClothesCategoryID:   NullInt64Value(item.CategoryID),
			ClothesCategoryName: NullStringValue(item.CategoryName),
			ServiceTypeID:       NullInt64Value(item.ServiceTypeID),
			ServiceTypeName:     NullStringValue(item.ServiceTypeName),
			CustomServiceName:   NullStringValue(item.CustomServiceName),
			Price:               item.Price,
			Notes:               NullStringValue(item.Notes),
			Status:              item.Status,
		})
	}

	combinedOrders := make([]*ListOrdersOutput, len(orders))

	// each order, get hashmap order items by orderID
	for i, order := range orders {
		var deadline *time.Time
		if order.Deadline.Valid {
			t := order.Deadline.Time
			deadline = &t
		}

		combinedOrders[i] = &ListOrdersOutput{
			OrderID:      order.ID,
			Name:         order.Name,
			Deadline:     deadline,
			CustomerName: order.CustomerName,
			CustomerId:   order.CustomerID,
			Status:       order.Status,
			Items:        orderMap[order.ID],
		}
	}

	return combinedOrders

}

type DetailOrderItemOutput struct {
	ClothesFor          string
	ClothesCategoryID   int64
	ClothesCategoryName string
	ServiceTypeID       int64
	ServiceTypeName     string
	CustomServiceName   string
	Price               int64
	Notes               string
	Status              string
}
type DetailOrderOutput struct {
	OrderID      int64
	Name         string
	Deadline     *time.Time
	CustomerName string
	CustomerId   int64
	Status       string
	Items        []DetailOrderItemOutput
}

func (o *orderService) GetOrderDetailsByOrderID(ctx context.Context, userId int64, orderID int64) (*DetailOrderOutput, error) {

	// get order
	order, err := o.store.GetOrderDetailsByOrderID(ctx, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &DetailOrderOutput{}, errors.New("order not found")
		}
		return &DetailOrderOutput{}, err
	}

	// get order_items
	orderItems, err := o.store.ListOrderItemsByOrderID(ctx, db.ListOrderItemsByOrderIDParams{
		OrderID: orderID,
		Limit:   100,
		Offset:  0,
	})
	if err != nil {
		return &DetailOrderOutput{}, err
	}

	if order.UserID != userId {
		return &DetailOrderOutput{}, errors.New("order does not belong to user")
	}

	var deadline *time.Time
	if order.Deadline.Valid {
		t := order.Deadline.Time
		deadline = &t
	}

	var orderItemsOutput []DetailOrderItemOutput
	for _, item := range orderItems {
		orderItemsOutput = append(orderItemsOutput, DetailOrderItemOutput{
			ClothesFor:          item.ClothesFor,
			ClothesCategoryID:   NullInt64Value(item.CategoryID),
			ClothesCategoryName: NullStringValue(item.CategoryName),
			ServiceTypeID:       NullInt64Value(item.ServiceTypeID),
			ServiceTypeName:     NullStringValue(item.ServiceTypeName),
			CustomServiceName:   NullStringValue(item.CustomServiceName),
			Price:               item.Price,
			Notes:               NullStringValue(item.Notes),
			Status:              item.Status,
		})
	}

	return &DetailOrderOutput{
		OrderID:      order.ID,
		Name:         order.Name,
		Deadline:     deadline,
		CustomerName: order.CustomerName,
		CustomerId:   order.CustomerID,
		Status:       order.Status,
		Items:        orderItemsOutput,
	}, nil
}
