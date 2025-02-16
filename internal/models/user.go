package models

type UserStatus string

const (
	Active   UserStatus = "ACTIVE"
	Inactive UserStatus = "INACTIVE"
)

type Address struct {
	CEP          string `json:"cep"`
	Address      string `json:"address"`
	Neighborhood string `json:"neighborhood"`
	State        string `json:"state"`
	City         string `json:"city"`
}

type User struct {
	ID        string     `json:"id" bson:"_id,omitempty"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	BirthDate string     `json:"birthDate"`
	Phone     string     `json:"phone"`
	Gender    string     `json:"gender"`
	CPF       string     `json:"cpf"`
	Address   Address    `json:"address"`
	Password  string     `json:"password,omitempty"`
	Status    UserStatus `json:"status"`
}
