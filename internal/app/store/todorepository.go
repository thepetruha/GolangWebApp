package store

import (
	"errors"
	"fmt"
	"webapp/internal/app/model"
)

type TodoRepository struct {
	store *Store
}

func (r *TodoRepository) CreateRow(t *model.ToDo) error {
	var row = r.store.db.QueryRow(
		`INSERT INTO todo (customer_id, text_todo, date_todo) VALUES ($1, $2, $3) RETURNING id_td`,
		t.CustomerID,
		t.Text,
		t.Date,
	)
	if row == nil {
		return errors.New("Error querying")
	}

	if err := row.Scan(&t.ID); err != nil {
		fmt.Println("Error querying insert user")
		return err
	}

	return nil
}

func (t *TodoRepository) FindByUserId(customer_id int) ([]*model.ToDo, error) {
	rows, err := t.store.db.Query("SELECT id_td, text_todo, date_todo FROM todo WHERE customer_id = $1", customer_id)
	if err != nil {
		return nil, errors.New("Error query todo list")
	}

	defer rows.Close()

	list := []*model.ToDo{}

	for rows.Next() {
		p := &model.ToDo{}

		err := rows.Scan(&p.ID, &p.Text, &p.Date)
		if err != nil {
			fmt.Println(err)
			continue
		}
		list = append(list, p)
	}
	for _, p := range list {
		fmt.Println(p.ID, p.Text, p.Date)
	}

	return list, nil
}
