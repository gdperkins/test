package gomap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type employee struct {
	FirstName string
	LastName  string
	Salary    int
}

type employeeViewModel struct {
	FirstName string
	LastName  string
	FullName  string
}

func TestCanMap(t *testing.T) {
	t.Log("Test can map with no configuration")

	source := employee{"John", "Doe", 1000}
	destination := employeeViewModel{}

	gm := New()
	gm.Map(source, &destination)

	assert.Equal(t, destination.FirstName, source.FirstName, "FirstName should be equal")
	assert.Equal(t, destination.LastName, source.LastName, "LastName should be equal")
	assert.Empty(t, destination.FullName)
}

func TestCanIgnoreField(t *testing.T) {
	t.Log("Ignoring FirstName (expecting: empty string)")

	source := employee{"John", "Doe", 1000}
	destination := employeeViewModel{}

	gm := New()
	gm.Add(source, destination, map[string]FieldConfig{
		"FirstName": {
			Ignore: true,
		},
	})
	gm.Map(source, &destination)
	assert.Empty(t, destination.FirstName, "FirstName should be empty")
}

func TestCanChangeFieldSource(t *testing.T) {
	t.Log("Changing FirstName source to LastName (expecting: Doe)")
	gm := New()

	source := employee{"John", "Doe", 1000}
	destination := employeeViewModel{}

	gm.Add(source, destination, map[string]FieldConfig{
		"FirstName": {
			Source: "LastName",
		},
	})

	gm.Map(source, &destination)
	assert.Equal(t, destination.FirstName, "Doe", "FirstName should be Doe")
	assert.Equal(t, destination.LastName, "Doe", "LastName should be Doe")
}

func TestMapTypeAcceptance(t *testing.T) {
	t.Log("Confirming source and destination types")

	source := employee{"John", "Doe", 1000}
	destination := employeeViewModel{}

	gm := New()

	gm.Map(&source, &destination)
	assert.Equal(t, destination.FirstName, "John", "FirstName should be John")
	gm.Map(source, &destination)
	assert.Equal(t, destination.FirstName, "John", "FirstName should be John")
	assert.Error(t, gm.Map(source, destination), "Should not be able to pass destination as value type")
	assert.Error(t, gm.Map(100, &destination), "Source should be of type struct")
}

func TestDuplicateMapConfigsIgnored(t *testing.T) {
	t.Log("Confirming duplicate configs are ignored")
	gm := New()

	source := employee{"John", "Doe", 1000}
	destination := employeeViewModel{}

	gm.Add(source, destination, map[string]FieldConfig{
		"FirstName": {
			Source: "LastName",
		},
	})

	gm.Add(source, destination, map[string]FieldConfig{
		"FirstName": {
			Source: "LastName",
		},
	})

	assert.Equal(t, 1, len(gm.mappingConfig), "The second map should have been ignored")
}
