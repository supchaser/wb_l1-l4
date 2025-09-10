package main

import (
	"crypto/rand"
	"fmt"
)

type GenderType string

const (
	MaleType   GenderType = "man"
	FemaleType GenderType = "woman"
)

type Human struct {
	FirstName string
	LastName  string
	Age       string
	Gender    GenderType
}

func CreateHuman(firstName, lastName, age string, gender GenderType) *Human {
	return &Human{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
		Gender:    gender,
	}
}

func (h *Human) GetFullName() string {
	return fmt.Sprintf("The full name this person is %s %s", h.FirstName, h.LastName)
}

func (h *Human) Speak() string {
	return rand.Text()
}

type Action struct {
	Human
}

func main() {
	firstName, lastName, age := "Valentin", "Stremin", "21"
	genderType := MaleType
	human := CreateHuman(firstName, lastName, age, genderType)

	fullNameHumanStruct := human.GetFullName()
	fmt.Println(fullNameHumanStruct)

	action := &Action{Human: *human}
	fullNameActionStruct := action.GetFullName()
	fmt.Println(fullNameActionStruct)

	fmt.Println(fullNameHumanStruct == fullNameActionStruct)

	for range 10 {
		fmt.Println(action.Speak())
	}
}
