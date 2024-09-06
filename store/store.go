package store

import (
	"database/sql"
	"errors"
)

var ErrorNotFound = errors.New("Record Not Found")

type Storage struct {
	db *sql.DB
}

type Store interface {
	//CreateCar(car *types.Car) (*types.Car, error)
	//GetCars() ([]types.Car, error)
	//DeleteCar(id string) error
	//FindCarsByNameMakeOrBrand(search string) ([]types.Car, error)
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) DeleteCar(id string) error {
	result, err := s.db.Exec("DELETE FROM cars WHERE id = ?", id)
    if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0{
        return ErrorNotFound
    }
	return err
}

//func (s *Storage) CreateCar(c *types.Car) (*types.Car, error) {
//	row, err := s.db.Exec("INSERT INTO cars (brand, make, model, year, imageURL) VALUES (?, ?, ?, ?, ?)", c.Brand, c.Make, c.Model, c.Year, c.ImageURL)
//	if err != nil {
//		return nil, err
//	}
//
//	id, err := row.LastInsertId()
//	if err != nil {
//		return nil, err
//	}
//	c.ID = int(id)
//
//	return c, nil
//}
//
//func (s *Storage) GetCars() ([]types.Car, error) {
//	rows, err := s.db.Query("SELECT * FROM cars")
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var cars []types.Car
//	for rows.Next() {
//		car, err := scanCar(rows)
//		if err != nil {
//			return nil, err
//		}
//		cars = append(cars, car)
//	}
//
//	return cars, nil
//}
