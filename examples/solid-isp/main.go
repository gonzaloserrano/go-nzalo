package main

import "github.com/gonzaloserrano/go-nzalo/examples/solid-isp/animal"

// Lets say we want to write a function to listen the sounds of some animals.
// The animal package provides us with a cow and cat, but we want a dog too.
// Because the animal.Animal interface is too broad with Speak, Breadth, Eat
// and Poo methods, we segregate it creating a speaker interface in this
// package (the client of the animal package) that just defines the Speak
// method of the animal.
type speaker interface {
	Speak()
}

// A dog just speaks, we don't need anything else.
type dog struct{}

func (d dog) Speak() {
	println("woof woof")
}

// listen implements the main functionality we were asked for: making animals speak.
func speak(speakers []speaker) {
	for _, speaker := range speakers {
		speaker.Speak()
	}
}

func main() {
	cow := animal.Cow{}
	cat := animal.Cat{}
	dog := dog{}

	// this works because cow and cat implement Speak. In this pckage we don't
	// care if they also implement other interfaces like animal.Animal.
	speak([]speaker{cow, cat, dog})
}
