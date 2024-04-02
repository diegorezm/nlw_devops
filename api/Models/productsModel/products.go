package productsmodel

import "time"

type ProductEntity struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Product struct {
	Id        int       `json:"id"`
	EntityId  int       `json:"entity_id"`
	Sold      bool      `json:"sold"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProductEntity(id int, name string, createdAt, updatedAt time.Time) *ProductEntity {
	return &ProductEntity{
		Id:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewProduct(id, entityId int, sold bool, createdAt, updatedAt time.Time) *Product {
	return &Product{
		Id:        id,
		EntityId:  entityId,
		Sold:      sold,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
