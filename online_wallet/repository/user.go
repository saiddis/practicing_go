package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/saiddis/practicing_go/online_wallet/domain"
	"github.com/saiddis/practicing_go/online_wallet/postgres"
)

type UserRepository struct {
	*postgres.PG
}

func NewUserRepository(db *postgres.PG) *UserRepository {
	return &UserRepository{
		PG: db,
	}
}

func (us UserRepository) FindUserByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	row := us.DB.QueryRow(ctx, "SELECT * FROM users WHERE id = $1;", id)

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Balance)
	if err != nil {
		log.Printf("Error retrieving user by id: %v", err)
	}
	return &user, nil
}

func (us *UserRepository) CreateUser(ctx context.Context, user domain.User) error {
	query := `INSERT INTO users (name, email, created_at, updated_at, balance) VALUES (
	@name,
	@email,
	@createdAt,
	@updatedAt,
	@balance
	);`
	args := pgx.NamedArgs{
		"name":      user.Name,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
		"balance":   user.Balance,
	}

	tag, err := us.DB.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("error inserting row: %v", err)
	}
	log.Printf("Create user: %v", tag)
	return nil
}

func (us *UserRepository) UpdateUser(ctx context.Context, id int, upd domain.UserUpdate) error {
	query := `UPDATE users 
	SET name = @name, email = @email, updated_at = @updatedAt
	WHERE id = @id;
	`
	args := pgx.NamedArgs{
		"id":        id,
		"name":      upd.Name,
		"email":     upd.Email,
		"updatedAt": time.Now().UTC(),
	}

	tag, err := us.DB.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("error updating record: %v", err)
	}
	log.Printf("Update user with id %d: %v", id, tag)
	return nil
}

func (us *UserRepository) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users
	WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}
	tag, err := us.DB.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("error deleting record: %v", err)
	}
	log.Printf("Delete user with id %d: %v", id, tag)
	return nil
}

func (us *UserRepository) Transfer(ctx context.Context, srcID, dstID int, amount float32) error {
	conn, err := us.DB.Acquire(ctx)
	if err != nil {
		return domain.Errorf(domain.EINTERNAL, "error acquiring connection: %v", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Errorf(domain.EINTERNAL, "error beginning transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {

			tx.Commit(ctx)
		}
	}()

	err = us.Withdraw(ctx, srcID, amount)
	if err != nil {
		return err
	}

	err = us.Credit(ctx, dstID, amount)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserRepository) Credit(ctx context.Context, id int, amount float32) error {
	user, err := us.FindUserByID(ctx, id)
	if err != nil {
		return domain.Errorf(domain.EINVALID, "error retrieving user with id %d: %v", id, err)
	}
	newBalance := user.Balance + amount
	if newBalance > domain.UnauthorizedBalanceLimit {
		return domain.Errorf(domain.EREACHEDLIMIT, "error adding up user's balance with id %d", user.ID)
	}
	query := `UPDATE users 
	SET balance = @balance, updated_at = @updatedAt
	WHERE id = @id;
	`
	args := pgx.NamedArgs{
		"id":        id,
		"balance":   newBalance,
		"updatedAt": time.Now().UTC(),
	}

	tag, err := us.DB.Exec(ctx, query, args)
	if err != nil {
		return domain.Errorf(domain.EINTERNAL, "error crediting user's balance: %v", err)
	}
	log.Printf("Credit to user with id %d: %v", id, tag)
	return nil
}

func (us *UserRepository) Withdraw(ctx context.Context, id int, amount float32) error {
	user, err := us.FindUserByID(ctx, id)
	if err != nil {
		return domain.Errorf(domain.EINVALID, "error retrieving user with id %d: %v", id, err)
	}
	newBalance := user.Balance - amount
	if newBalance < 0 {
		return domain.Errorf(domain.ENOTENOUGH, "error withdrawing from user with id %d", user.ID)
	}
	query := `UPDATE users 
	SET balance = @balance, updated_at = @updatedAt
	WHERE id = @id;
	`
	args := pgx.NamedArgs{
		"id":        id,
		"balance":   newBalance,
		"updatedAt": time.Now().UTC(),
	}

	tag, err := us.DB.Exec(ctx, query, args)
	if err != nil {
		return domain.Errorf(domain.EINTERNAL, "error crediting user's balance: %v", err)
	}
	log.Printf("Withdraw from user with id %d: %v", id, tag)
	return nil
}
