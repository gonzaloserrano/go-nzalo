package animal

type Animal interface {
	Speak()
	Breath()
	Eat()
	Poo()
}

type Cow struct{}

func (c Cow) Speak() {
	println("muuuuh")
}
func (c Cow) Breath() {}
func (c Cow) Eat()    {}
func (c Cow) Poo()    {}

type Cat struct{}

func (c Cat) Speak() {
	println("meoow")
}
func (c Cat) Breath() {}
func (c Cat) Eat()    {}
func (c Cat) Poo()    {}
