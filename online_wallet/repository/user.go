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

func (us UserRepository) FindUserByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	err := us.DB.QueryRow(ctx, "select * from users where id = $1;", id).Scan(&user)
	if err != nil {
		log.Printf("QueryRow failed: %v", err)
	}
	return &user, nil
}

func (us *UserRepository) CreateUser(ctx context.Context, user domain.User) error {
	query := `INSERT INTO users (name, email) VALUES (
	@name,
	@email,
	@createdAt,
	@updatedAt,
	@balace
	);`
	args := pgx.NamedArgs{
		"name":      user.Name,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
		"balance":   user.Balance,
	}

	_, err := us.DB.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("error inserting row: %v", err)
	}
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

	_, err := us.DB.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("error updating record: %v", err)
	}
	return nil
}

func (us *UserRepository) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users
	WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}
	_, err := us.DB.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("error deleting record: %v", err)
	}
	return nil
}

func (us *UserRepository) Transfer(ctx context.Context, srcID, dstID int, amount float64) error {
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

func (us *UserRepository) Credit(ctx context.Context, id int, amount float64) error {
	user, err := us.FindUserByID(ctx, id)
	if err != nil {
		return domain.Errorf(domain.EINVALID, "error retrieving user with id:%d: %v", id, err)
	}
	newBalance := user.Balance + amount
	if newBalance > domain.UnauthorizedBalanceLimit {
		return domain.Errorf(domain.EREACHEDLIMIT, "error adding up user's balance with id:%d", user.ID)
	}
	query := `UPDATE users 
	SET balance = balance, updated_at = @updatedAt
	WHERE id = @id;
	`
	args := pgx.NamedArgs{
		"id":        id,
		"balance":   newBalance,
		"updatedAt": time.Now().UTC(),
	}

	_, err = us.DB.Exec(ctx, query, args)
	if err != nil {
		return domain.Errorf(domain.EINTERNAL, "error crediting user's balance: %v", err)
	}
	return nil
}

func (us *UserRepository) Withdraw(ctx context.Context, id int, amount float64) error {
	user, err := us.FindUserByID(ctx, id)
	if err != nil {
		return domain.Errorf(domain.EINVALID, "error retrieving user with id:%d: %v", id, err)
	}
	newBalance := user.Balance - amount
	if newBalance < 0 {
		return domain.Errorf(domain.ENOTENOUGH, "error withdrawing from user with id:%d", user.ID)
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

	_, err = us.DB.Exec(ctx, query, args)
	if err != nil {
		return domain.Errorf(domain.EINTERNAL, "error crediting user's balance: %v", err)
	}
	return nil
}
