package abi

import (
	"reflect"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

type Address struct {
	Street string
	City   string
	State  string
}

type Employee struct {
	Person  Person
	Address Address
}

func TestDeepClone(t *testing.T) {
	// Test case 1: Clone a struct with non-struct fields
	p1 := Person{Name: "Alice", Age: 30}
	p2 := deepClone(p1).(Person)
	if !reflect.DeepEqual(p1, p2) {
		t.Errorf("Expected %v, but got %v", p1, p2)
	}

	// Test case 2: Clone a struct with nested structs
	e1 := Employee{
		Person:  Person{Name: "Bob", Age: 40},
		Address: Address{Street: "123 Main St", City: "Anytown", State: "CA"},
	}
	e2 := deepClone(e1).(Employee)
	if !reflect.DeepEqual(e1, e2) {
		t.Errorf("Expected %v, but got %v", e1, e2)
	}

	// Test case 3: Clone a non-struct value
	s1 := "hello"
	s2 := deepClone(s1).(string)
	if s1 != s2 {
		t.Errorf("Expected %v, but got %v", s1, s2)
	}
}
