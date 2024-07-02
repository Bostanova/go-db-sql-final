package main

import (
	"database/sql"
)

// ParcelStore реализует логику работы с данными о посылке в БД
type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// Add добавляет новую посылку в БД и возвращает
// id добавленной посылки или ошибку
func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) "+
		"VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get возвращает посылку с номером number или ошибку
func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}

	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return []Parcel{}, err
	}
	defer rows.Close()

	var res []Parcel

	for rows.Next() {
		p := Parcel{}

		err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return []Parcel{}, err
		}

		if err = rows.Err(); err != nil {
			return []Parcel{}, err
		}

		res = append(res, p)
	}
	if err = rows.Err(); err != nil {
		return []Parcel{}, err
	}
	return res, nil
}

// SetStatus обновляет статус посылки с номером number на status
func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

// SetAddress устанавливает новый адрес доставки, если статус посылки не "registered"
func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE NUMBER = :number and status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}

	return nil
}

// Delete удаляет из БД посылку с номером number или ошибку
func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("DELETE FROM parcel WHERE NUMBER = :number and status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}

	return nil
}
