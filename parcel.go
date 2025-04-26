package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

const (
	insertParcel   = "INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address,:created_at)"
	selectByNumber = "SELECT client, status, address, created_at FROM parcel WHERE number = :number"
	selectByClient = "SELECT number, client, status, address, created_at FROM parcel WHERE client = :client"
	updateStatus   = "UPDATE parcel SET status = :status WHERE number = :number"
	updateAddress  = "UPDATE parcel SET address = :address WHERE number = :number AND status = 'registered'"
	deleteParcel   = "DELETE FROM parcel WHERE number = :number AND status = 'registered'"
)

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec(insertParcel,
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt),
	)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow(selectByNumber, sql.Named("number", number))

	p := Parcel{}

	err := row.Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query(selectByClient, sql.Named("client", client))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res []Parcel

	for rows.Next() {
		var p Parcel
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec(updateStatus,
		sql.Named("status", status),
		sql.Named("number", number),
	)

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec(updateAddress,
		sql.Named("address", address),
		sql.Named("number", number),
	)

	return err
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec(deleteParcel, sql.Named("number", number))
	return err
}
