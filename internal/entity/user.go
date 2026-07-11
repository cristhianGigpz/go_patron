package entity

type User struct {
	ID       uint
	Name     string
	Email    string
	Age      int
	Password string
	Role     string
}

//No existe ningún import. Eso hace que el modelo sea independiente.
