package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	db "jahitin_be/database/repository"
	"time"

	"github.com/nyaruka/phonenumbers"
)

type OrderService interface {
	CreateUserOrder(ctx context.Context, user_id int64, o_data CreateUserOrderInput, c_data CustomerInput) (bool, error)
	ListUserOrders(ctx context.Context, filter ListOrdersFilter) ([]*ListOrdersOutput, error)
	GetOrderDetailsByOrderID(ctx context.Context, userId int64, orderID int64) (*DetailOrderOutput, error)
	ReplaceOrderItemsByOrderID(ctx context.Context, userId int64, orderItemData []OrderItemUpdateInput, orderID int64) error
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
	ClothesCategoryID   *int64
	ServiceTypeID       *int64
	CustomServiceName   *string
	Price               int64
	IsSaveCustomerNotes bool
}
type CreateUserOrderInput struct {
	Name     string
	Deadline *time.Time
	Items    []OrderItemInput
}

type CustomerInput struct {
	Name  string
	Phone string
}

func (o *orderService) CreateUserOrder(ctx context.Context, user_id int64, o_data CreateUserOrderInput, c_data CustomerInput) (bool, error) {
	err := o.store.ExecTx(ctx, func(q *db.Queries) error {

		var customerID, err = o.resolveCustomer(ctx, q, user_id, c_data)
		if err != nil {
			return err
		}

		fmt.Println("customer ID	", customerID)

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

			// if serviceTypeID is nil, cusom service name must be provided
			if item.ServiceTypeID == nil && (item.CustomServiceName == nil || *item.CustomServiceName == "") {
				return errors.New("custom service name must be provided if service type ID is nil")
			}

			oiParam := db.CreateOrderItemParams{
				OrderID:           order.ID,
				CategoryID:        nullInt64(item.ClothesCategoryID),
				ServiceTypeID:     nullInt64(item.ServiceTypeID),
				ClothesFor:        item.ClothesFor,
				CustomServiceName: nullString(item.CustomServiceName),
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

// get customer by phone number, if not found, create new customer
func (o *orderService) resolveCustomer(ctx context.Context, q db.Querier, user_id int64, customerData CustomerInput) (customerID int64, err error) {

	formattedPhone, err := phonenumbers.Parse(customerData.Phone, "ID")
	if err != nil {
		return 0, err
	}
	formattedPhoneE164 := phonenumbers.Format(formattedPhone, phonenumbers.E164)

	customer, _ := o.store.GetCustomerByUserIDAndPhone(ctx, db.GetCustomerByUserIDAndPhoneParams{
		UserID:         user_id,
		FormattedPhone: sql.NullString{String: formattedPhoneE164, Valid: true},
	})

	if customer.ID != 0 {
		return customer.ID, nil

	} else {
		// create customer
		cParam := db.CreateCustomerParams{
			Name:           customerData.Name,
			UserID:         user_id,
			FormattedPhone: sql.NullString{String: formattedPhoneE164, Valid: true},
			CountryCode:    sql.NullString{String: fmt.Sprintf("+%d", formattedPhone.GetCountryCode()), Valid: true},
			Phone:          sql.NullString{String: customerData.Phone, Valid: customerData.Phone != ""},
			Notes:          sql.NullString{String: "", Valid: false},
		}

		customer, err := q.CreateCustomer(ctx, cParam)
		if err != nil {
			return 0, err
		}

		return customer.ID, nil
	}

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
	TotalPrice   int64
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

	var totalPrice int64 = 0

	var orderItemsOutput []DetailOrderItemOutput
	for _, item := range orderItems {
		totalPrice += item.Price

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
		TotalPrice:   totalPrice,
	}, nil
}

type OrderItemUpdateInput struct {
	ClothesFor          string
	Notes               string
	ClothesCategoryID   *int64
	ServiceTypeID       *int64
	CustomServiceName   *string
	Price               int64
	IsSaveCustomerNotes bool
}

func (o *orderService) ReplaceOrderItemsByOrderID(ctx context.Context, userId int64, orderItemData []OrderItemUpdateInput, orderID int64) error {

	// get order
	existingOrder, err := o.store.GetOrderDetailsByOrderID(ctx, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("order not found")
		}
		return err
	}

	if existingOrder.UserID != userId {
		return errors.New("order does not belong to user")
	}

	return o.store.ExecTx(ctx, func(q *db.Queries) error {
		// Replace all existing order items atomically.
		if err := q.DeleteOrderItemByOrderID(ctx, orderID); err != nil {
			return err
		}

		// Insert items individually so nullable foreign keys remain SQL NULL.
		for _, item := range orderItemData {
			_, err := q.CreateOrderItem(ctx, db.CreateOrderItemParams{
				OrderID:           existingOrder.ID,
				CategoryID:        nullInt64(item.ClothesCategoryID),
				ServiceTypeID:     nullInt64(item.ServiceTypeID),
				ClothesFor:        item.ClothesFor,
				CustomServiceName: nullString(item.CustomServiceName),
				Notes:             sql.NullString{String: item.Notes, Valid: item.Notes != ""},
				Price:             item.Price,
			})
			if err != nil {
				return err
			}
		}

		// Save customer notes in the same transaction.
		for _, item := range orderItemData {
			if item.IsSaveCustomerNotes {
				if err := o.saveCustomerOrderNotes(ctx, q, item.Notes, existingOrder.CustomerID); err != nil {
					return err
				}
			}
		}

		return nil
	})

}
