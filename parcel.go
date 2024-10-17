package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// Подготовка SQL-запроса
	stmt, err := s.db.Prepare("INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Выполнение запроса
	result, err := stmt.Exec(p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}

	// Получение ID новой записи
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// Подготовка SQL-запроса
	stmt, err := s.db.Prepare("SELECT * FROM parcel WHERE number = ?")
	if err != nil {
		return Parcel{}, err
	}
	defer stmt.Close()

	// Выполнение запроса
	row := stmt.QueryRow(number)

	// Считывание данных из строки
	var p Parcel
	if err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// Подготовка SQL-запроса
	stmt, err := s.db.Prepare("SELECT * FROM parcel WHERE client = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Выполнение запроса
	rows, err := stmt.Query(client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Считывание данных из строк
	var parcels []Parcel
	for rows.Next() {
		var p Parcel
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}

	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// Подготовка SQL-запроса
	stmt, err := s.db.Prepare("UPDATE parcel SET status = ? WHERE number = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполнение запроса
	_, err = stmt.Exec(status, number)
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// Подготовка SQL-запроса
	stmt, err := s.db.Prepare("UPDATE parcel SET address = ? WHERE number = ? AND status = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполнение запроса
	_, err = stmt.Exec(address, number, ParcelStatusRegistered)
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// Подготовка SQL-запроса
	stmt, err := s.db.Prepare("DELETE FROM parcel WHERE number = ? AND status = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполнение запроса
	_, err = stmt.Exec(number, ParcelStatusRegistered)
	if err != nil {
		return err
	}

	return nil
}
