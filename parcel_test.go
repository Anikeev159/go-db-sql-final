package main

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestParcelStore_Add(t *testing.T) {
	db, err := sql.Open("sqlite3", "./tracker.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := Parcel{
		Client:    1,
		Status:    ParcelStatusRegistered,
		Address:   "Адрес 1",
		CreatedAt: "2023-10-26T12:00:00Z",
	}

	id, err := store.Add(parcel)
	if err != nil {
		t.Fatal(err)
	}

	if id == 0 {
		t.Error("ID посылки должен быть больше нуля")
	}
}

func TestParcelStore_GetByClient(t *testing.T) {
	db, err := sql.Open("sqlite3", "./tracker.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	client := 1

	parcels, err := store.GetByClient(client)
	if err != nil {
		t.Fatal(err)
	}

	if len(parcels) == 0 {
		t.Error("Должен быть хотя бы один результат")
	}

	for _, parcel := range parcels {
		if parcel.Client != client {
			t.Errorf("Неверный клиент для посылки: %d", parcel.Client)
		}
	}
}

func TestParcelStore_UpdateStatus(t *testing.T) {
	db, err := sql.Open("sqlite3", "./tracker.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcelNumber := 1
	newStatus := ParcelStatusSent

	err = store.UpdateStatus(parcelNumber, newStatus)
	if err != nil {
		t.Fatal(err)
	}

	// Проверяем, что статус изменился
	rows, err := db.Query("SELECT status FROM parcel WHERE number = ?", parcelNumber)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var status string
	if rows.Next() {
		if err := rows.Scan(&status); err != nil {
			t.Fatal(err)
		}
	}

	if status != newStatus {
		t.Errorf("Статус посылки не изменился: %s", status)
	}
}

func TestParcelStore_UpdateAddress(t *testing.T) {
	db, err := sql.Open("sqlite3", "./tracker.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcelNumber := 1
	newAddress := "Адрес 2"

	err = store.UpdateAddress(parcelNumber, newAddress)
	if err != nil {
		t.Fatal(err)
	}

	// Проверяем, что адрес изменился
	rows, err := db.Query("SELECT address FROM parcel WHERE number = ?", parcelNumber)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var address string
	if rows.Next() {
		if err := rows.Scan(&address); err != nil {
			t.Fatal(err)
		}
	}

	if address != newAddress {
		t.Errorf("Адрес посылки не изменился: %s", address)
	}
}

func TestParcelStore_Delete(t *testing.T) {
	db, err := sql.Open("sqlite3", "./tracker.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcelNumber := 1

	err = store.Delete(parcelNumber)
	if err != nil {
		t.Fatal(err)
	}

	// Проверяем, что посылка удалена
	rows, err := db.Query("SELECT * FROM parcel WHERE number = ?", parcelNumber)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		t.Errorf("Посылка не удалена: %d", parcelNumber)
	}
}
