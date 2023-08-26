package postge

import (
	"database/sql"
	"errors"
	"math/rand"
	"test_pizza/internal/lib/random"
	"test_pizza/internal/models"
	"test_pizza/internal/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(DBSN string) (*Storage, error) {
	db, err := sql.Open("postgres", DBSN)
	if err != nil {
		return &Storage{}, err
	}
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS public.orders
	(
		order_id character varying(15),
		done boolean DEFAULT false,
		PRIMARY KEY (order_id),
		CONSTRAINT orders_ids UNIQUE (order_id)
	);`)
	if err != nil {
		return &Storage{}, err
	}
	_, err = stmt.Exec()
	if err != nil {
		return &Storage{}, err
	}
	// Можно писать items в строку в рамках таблицы orders, но так веселее
	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS public.items
	(
		id bigserial,
		order_id character varying(15) NOT NULL,
		item integer,
		PRIMARY KEY (id)
	);`)
	if err != nil {
		return &Storage{}, err
	}
	_, err = stmt.Exec()
	if err != nil {
		return &Storage{}, err
	}
	

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) CreateOrder(items []int) (models.Order, error){
	var resp models.Order
	var id string
	isExist := true
	for isExist {
		var res string
		id = random.NewRandomString(3+rand.Intn(13))
		stmt, err := s.db.Prepare("SELECT order_id FROM orders WHERE order_id = $1")
		if err!= nil {
			return models.Order{}, err
		}
		err = stmt.QueryRow(id).Scan(&res)
		if errors.Is(err, sql.ErrNoRows) {
			isExist = false
		} else {
			return models.Order{}, err
		}
	}
	
	tx, err := s.db.Begin()
	if err!= nil {
		return models.Order{}, err
	}

	stmtOrder, err := s.db.Prepare(`
	INSERT INTO orders (order_id)
	VALUES ($1);
	`)
	if err!= nil {
		tx.Rollback()
		return models.Order{}, err
	}
	defer stmtOrder.Close()
	_, err = stmtOrder.Exec(id)
	if err!= nil {
		tx.Rollback()
		return models.Order{}, err
	}

	stmtItems, err := s.db.Prepare(`
	INSERT INTO items (order_id, item)
	VALUES ($1, $2)
	`)
	if err!= nil {
		tx.Rollback()
		return models.Order{}, err
	}
	for _, item := range items {
		_, err = stmtItems.Exec(id, item)
		if err!= nil {
			tx.Rollback()
			return models.Order{}, err
		}
	}
	defer stmtItems.Close()

	err = tx.Commit()
	if err!= nil {
		return models.Order{}, err
	}


	resp.Order_id = id
	resp.Items = items
	resp.Done = false
	return resp, nil
}
func (s *Storage) AddItems(order_id string, items []int) error{
	var done bool
	stmt, err := s.db.Prepare("SELECT done FROM orders WHERE order_id = $1")
	if err!= nil {
		return err
	}
	err = stmt.QueryRow(order_id).Scan(&done)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFoundOrder
		} else {
			return err
		}
	}	

	if done {
		return storage.ErrOrderIsFinished
	}

	stmtItems, err := s.db.Prepare(`
	INSERT INTO items (order_id, item)
	VALUES ($1, $2);`)
	if err!= nil {
		return err
	}
	for _, item := range items {
		_, err = stmtItems.Exec(order_id, item)
		if err!= nil {
			return err
		}
	}
	return nil
}
func (s *Storage) GetOrderById(order_id string) (models.Order, error){
	var resp models.Order

	stmt, err := s.db.Prepare("SELECT order_id, done FROM orders WHERE order_id = $1")
	if err!= nil {
		return models.Order{}, err
	}
	err = stmt.QueryRow(order_id).Scan(&resp.Order_id, &resp.Done)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Order{}, storage.ErrNotFoundOrder
		} else {
			return models.Order{}, err
		}
	}	
	stmt, err = s.db.Prepare("SELECT item FROM items WHERE order_id = $1")
	if err!= nil {
		return models.Order{}, err
	}
	rows, err := stmt.Query(order_id)
	if err!= nil {
		return models.Order{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var item int
		if err := rows.Scan(&item); err != nil {
			return models.Order{}, err
		}
		resp.Items = append(resp.Items, item)
	}


	return resp, nil
}
func (s *Storage) FinishOrder(order_id string) error {
	var done bool
	stmt, err := s.db.Prepare("SELECT done FROM orders WHERE order_id = $1;")
	if err!= nil {
		return err
	}
	err = stmt.QueryRow(order_id).Scan(&done)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFoundOrder
		} else {
			return err
		}
	}
	if done {
		return storage.ErrOrderIsFinished
	}	
	stmt, err = s.db.Prepare(`UPDATE orders SET done=true WHERE order_id = $1;`)
	if err!= nil {
		return err
	}
	if _, err = stmt.Exec(order_id); err != nil {
		return err
	}
	return nil
}
func (s *Storage) GetOrdersByStatus(done int)  ([]models.Order, error) {
	var resp []models.Order
	var sql string
	sql = `SELECT order_id, done FROM orders`
	if done == 1 {
		sql += " WHERE done=true"
	}
	if done==0 {
		sql += " WHERE done=false"
	}

	stmt, err := s.db.Prepare(sql)
	if err!= nil {
		return []models.Order{}, err
	}
	rows, err := stmt.Query()
	if err!= nil {
		return []models.Order{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.Order_id, &order.Done); err != nil {
			return []models.Order{}, err
		}
		stmtItems, err := s.db.Prepare("SELECT item FROM items WHERE order_id = $1")
		if err!= nil {
			return []models.Order{}, err
		}
		items, err := stmtItems.Query(order.Order_id)
		if err!= nil {
			return []models.Order{}, err
		}
		for items.Next() {
			var item int
			if err := items.Scan(&item); err != nil {
				return []models.Order{}, err
			}
			order.Items = append(order.Items, item)
		}
		resp = append(resp, order)

	}


	return resp, nil
}