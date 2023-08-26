package ram

import (
	"math/rand"
	"test_pizza/internal/lib/random"
	models "test_pizza/internal/models"
	"test_pizza/internal/storage"
)



type Storage struct {
	orders []models.Order
}

func New() (*Storage) {
	return &Storage{
		orders: []models.Order{},
	}
}

func (s *Storage) CreateOrder(items []int) (models.Order, error) {
	var id string
	if (len(items)==0) {
		return models.Order{}, storage.ErrEmptyItems
	}
	for {
		id = random.NewRandomString(3+rand.Intn(13))
		isExist := false
		for _, o := range s.orders{
			if o.Order_id == id {
				isExist = true
			}
		}
		if !isExist {
			break
		}
	}
	order := models.Order{
		Order_id: id,
		Items: items,
		Done: false,
	}
	s.orders = append(s.orders, order)
	return order, nil
}

func (s *Storage) AddItems(order_id string, items []int) error {	
	for i, o := range s.orders{
		if o.Order_id == order_id {
			if o.Done {
				return storage.ErrOrderIsFinished
			}
			s.orders[i].Items = append(s.orders[i].Items, items...)
			return nil
		}
	}
	return storage.ErrNotFoundOrder
}

func (s *Storage) GetOrderById(order_id string) (models.Order, error) {
	for _, o := range s.orders{
		if o.Order_id == order_id {
			return o, nil
		}
	}
	return models.Order{}, storage.ErrNotFoundOrder
}

func (s *Storage) FinishOrder(order_id string) error {
	for i, o := range s.orders{
		if o.Order_id == order_id {
			if o.Done {
				return storage.ErrOrderIsFinished
			}
			s.orders[i].Done = true
			return nil
		}
	}
	return storage.ErrNotFoundOrder
}

func (s *Storage) GetOrdersByStatus(done int)  ([]models.Order, error){
	// 0 filter when Done = false
	// 1 filter when Done = true
	// -1 without filtering
	if done == -1 {
		return s.orders, nil
	}

	var doneStatus bool
	if done == 1 {
		doneStatus = true
	} else { doneStatus= false}

	var out []models.Order
	for _,o := range s.orders {
		if o.Done == doneStatus {
			out = append(out, o)
		}
	}
	return out, nil
}

	


