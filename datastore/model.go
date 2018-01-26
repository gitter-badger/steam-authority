package datastore

import (
	"time"

	"cloud.google.com/go/datastore"
)

type Model struct {
	Key       *datastore.Key `json:"-" datastore:"-"`
	ParentKey *datastore.Key `json:"-" datastore:"-"`
	CreatedAt time.Time      `json:"-" datastore:"created_at"`
	UpdatedAt time.Time      `json:"-" datastore:"updated_at"`
}

func (kind *Model) GetKey() *datastore.Key {
	return kind.Key
}

func (kind *Model) GetKeyAsString() string {
	return kind.GetKey().Encode()
}

func (kind *Model) Save() {

	kind.UpdatedAt = time.Now()

	//if kind.CreatedAt == nil {
	//	kind.CreatedAt = time.Now()
	//}
}

func (kind *Model) Delete() {
}

func (kind *Model) Reflect() {
	//val := reflect.ValueOf(kind).Elem()
	//
	//for i := 0; i < val.NumField(); i++ {
	//	valueField := val.Field(i)
	//	typeField := val.Type().Field(i)
	//	tag := typeField.Tag
	//
	//	fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("model"))
	//}
}
