package model

import (
	"database/sql"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	All() []User
	Save(*User) error
	Find(int) *User
	Login(string) *User
	FindRole(int) []User
	Update(*User) error
	Delete(user *User) error
}

type UserStoreMySQL struct {
	DB *sql.DB
}

type User struct {
	ID       int
	Name     string
	Address  string
	Telp     int
	Email    string
	Password string
	Role     string
	Token    string
}

func NewUserMySQL() UserStore {
	dsn := os.Getenv("DATABASE_USER") + ":" + os.Getenv("DATABASE_PASSWORD") + "@tcp(" + os.Getenv("DATABASE_HOST") + ")/" + os.Getenv("DATABASE_NAME") + "?parseTime=true&clientFoundRows=true"
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}
	return &UserStoreMySQL{DB: db}
}

func (store *UserStoreMySQL) All() []User {
	users := []User{}
	rows, err := store.DB.Query("SELECT * FROM user")

	if err != nil {
		return users
	}
	user := User{}
	for rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Address, &user.Telp, &user.Email, &user.Password, &user.Role, &user.Token)
		users = append(users, user)
	}
	return users
}

func CreateUser(name, address string, telp int, email, password, role, token string) (*User, error) {
	return &User{
		Name:     name,
		Address:  address,
		Telp:     telp,
		Email:    email,
		Password: password,
		Role:     role,
		Token:    token,
	}, nil
}

func (store *UserStoreMySQL) Save(user *User) error {
	result, err := store.DB.Exec(`
		INSERT INTO user(name, address, telp, email, password, role, token) VALUES(?,?,?,?,?,?,?)`,
		user.Name,
		user.Address,
		user.Telp,
		user.Email,
		user.Password,
		user.Role,
		user.Token,
	)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	_, err = result.RowsAffected()

	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()

	if err != nil {
		return err
	}
	user.ID = int(lastID)

	return nil
}

func (store *UserStoreMySQL) Find(id int) *User {
	user := User{}

	err := store.DB.
		QueryRow(`SELECT * FROM user WHERE id=?`, id).
		Scan(
			&user.ID,
			&user.Name,
			&user.Address,
			&user.Telp,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.Token,
		)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &user
}

func (store *UserStoreMySQL) FindRole(role int) []User {
	users := []User{}
	rows, err := store.DB.Query("SELECT * FROM user WHERE role=?", role)

	if err != nil {
		return users
	}
	user := User{}
	for rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Address, &user.Telp, &user.Email, &user.Password, &user.Role, &user.Token)
		users = append(users, user)
	}
	return users
}

func (store *UserStoreMySQL) Update(user *User) error {
	result, err := store.DB.Exec(`
		UPDATE user SET name = ?, address = ?, telp = ?, email = ?, password = ? WHERE id =?`,
		user.Name,
		user.Address,
		user.Telp,
		user.Email,
		user.Password,
		user.ID,
	)

	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil

}

func (store *UserStoreMySQL) Delete(user *User) error {
	result, err := store.DB.Exec(`
	DELETE FROM user WHERE id = ?`,
		user.ID,
	)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return nil
	}
	return nil
}

func (store *UserStoreMySQL) Login(email string) *User {
	user := User{}

	err := store.DB.
		QueryRow(`SELECT * FROM user WHERE email=?`, email).
		Scan(
			&user.ID,
			&user.Name,
			&user.Address,
			&user.Telp,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.Token,
		)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &user
}

func Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	}

	return string(hashed), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
