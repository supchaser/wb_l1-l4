package main

import (
	"crypto/rand"
	"fmt"
)

// Создал NewType
type GenderType string

// Задал константы
const (
	MaleType   GenderType = "man"
	FemaleType GenderType = "woman"
)

// Создал структуру Human с несколькими полями
type Human struct {
	FirstName string
	LastName  string
	Age       string
	Gender    GenderType
}

// Конструктор для структуры Human
func CreateHuman(firstName, lastName, age string, gender GenderType) *Human {
	return &Human{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
		Gender:    gender,
	}
}

// Метод получения полного имени для Human
func (h *Human) GetFullName() string {
	return fmt.Sprintf("The full name this person is %s %s", h.FirstName, h.LastName)
}

// Метод структуры Human, который генерирует рандомную строчку
func (h *Human) Speak() string {
	return rand.Text()
}

// Создаем структуру Action, в которую встраиваем стр-ру Human
type Action struct {
	Human
}

func main() {
	firstName, lastName, age := "Valentin", "Stremin", "21"
	genderType := MaleType

	// Инициализируем Human с помощью конструктора
	human := CreateHuman(firstName, lastName, age, genderType)

	// Получаем полное имя для струтктуры Human и распичатываем
	fullNameHumanStruct := human.GetFullName()
	fmt.Println(fullNameHumanStruct)

	// Создаем стр-ру Action
	action := &Action{Human: *human}

	// Т.к. мы встроили в Action стр-ру Human, то теперь имеем доступ ко всем методам стр-ры Human
	fullNameActionStruct := action.GetFullName()
	fmt.Println(fullNameActionStruct)

	fmt.Println(fullNameHumanStruct == fullNameActionStruct)

	for range 10 {
		fmt.Println(action.Speak())
	}
}
