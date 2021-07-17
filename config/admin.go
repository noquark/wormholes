package config

type Admin struct {
	Email    string
	Password string
}

func DefaultAdmin() Admin {
	return Admin{
		Email:    "admin@example.com",
		Password: "wormholes",
	}
}
