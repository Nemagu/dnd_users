package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Nemagu/dnd/internal/application"
	"github.com/Nemagu/dnd/internal/domain"
	"github.com/Nemagu/dnd/internal/domain/duser"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	userTable        = "\"user\""
	userVersionTable = "user_version"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		pool: pool,
	}
}

func (r *PostgresUserRepository) NextID(ctx context.Context) uuid.UUID {
	return uuid.New()
}

func (r *PostgresUserRepository) IDExists(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	var exists bool
	query := fmt.Sprintf(
		"SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)",
		userTable,
	)
	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return exists, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	return exists, nil
}

func (r *PostgresUserRepository) EmailExists(
	ctx context.Context,
	email domain.Email,
) (bool, error) {
	var exists bool
	query := fmt.Sprintf(
		"SELECT EXISTS(SELECT 1 FROM %s WHERE email = $1)",
		userTable,
	)
	err := r.pool.QueryRow(ctx, query, email.String()).Scan(&exists)
	if err != nil {
		return exists, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	return exists, nil
}

func (r *PostgresUserRepository) All(
	ctx context.Context,
	limit, offset int,
) ([]*duser.User, error) {
	query := fmt.Sprintf(
		`SELECT
			id,
			email,
			state,
			status,
			password_hash,
			version
		FROM %s`,
		userTable,
	)
	values := make([]any, 0, 2)
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", len(values)+1)
		values = append(values, limit)
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", len(values)+1)
		values = append(values, offset)
	}
	query += " ORDER BY email"
	rows, err := r.pool.Query(ctx, query, values...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return make([]*duser.User, 0), nil
		}
		return nil, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	defer rows.Close()
	cap := 30
	if limit != 0 {
		cap = limit
	}
	result := make([]*duser.User, 0, cap)
	for rows.Next() {
		u, err := buildUserFromRow(rows, "")
		if err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, nil
}

func (r *PostgresUserRepository) ByID(
	ctx context.Context,
	id uuid.UUID,
) (*duser.User, error) {
	query := fmt.Sprintf(
		`SELECT 
			id,
			email,
			state,
			status,
			password_hash,
			version 
		FROM %s
		WHERE id = $1
		ORDER BY email`,
		userTable,
	)
	row := r.pool.QueryRow(ctx, query, id)
	return buildUserFromRow(
		row,
		fmt.Sprintf("пользователя с id %s не существует", id),
	)
}

func (r *PostgresUserRepository) ByEmail(
	ctx context.Context,
	email domain.Email,
) (*duser.User, error) {
	query := fmt.Sprintf(
		`SELECT 
		    id,
			email,
			state,
			status,
			password_hash,
			version 
		FROM %s
		WHERE email = $1
		ORDER BY email`,
		userTable,
	)
	row := r.pool.QueryRow(ctx, query, email.String())
	return buildUserFromRow(
		row,
		fmt.Sprintf("пользователя с email %s не существует", email),
	)
}

func (r *PostgresUserRepository) Filter(
	ctx context.Context,
	searchByEmail string,
	filterByState []duser.UserState,
	filterByStatus []duser.UserStatus,
	limit, offset int,
) ([]*duser.User, error) {
	query := fmt.Sprintf(
		`SELECT
		    id,
			email,
			state,
			status,
			password_hash,
			version 
		FROM %s`,
		userTable,
	)
	conditions := make([]string, 0, 3)
	values := make([]any, 0, 1+len(filterByState)+len(filterByStatus)+2)
	if searchByEmail != "" {
		condition := fmt.Sprintf("LOWER(email) LIKE $%d", len(values)+1)
		conditions = append(conditions, condition)
		values = append(values, "%"+strings.ToLower(searchByEmail)+"%")
	}
	if len(filterByState) != 0 {
		temp := make([]string, 0, len(filterByState))
		for _, state := range filterByState {
			temp = append(temp, fmt.Sprintf("$%d", len(values)+1))
			values = append(values, state.String())
		}
		condition := "state IN (" + strings.Join(temp, ",") + ")"
		conditions = append(conditions, condition)
	}
	if len(filterByStatus) != 0 {
		temp := make([]string, 0, len(filterByStatus))
		for _, status := range filterByStatus {
			temp = append(temp, fmt.Sprintf("$%d", len(values)+1))
			values = append(values, status.String())
		}
		condition := "status IN (" + strings.Join(temp, ",") + ")"
		conditions = append(conditions, condition)
	}
	if len(values) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", len(values)+1)
		values = append(values, limit)
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", len(values)+1)
		values = append(values, offset)
	}
	query += " ORDER BY email"
	rows, err := r.pool.Query(ctx, query, values...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return make([]*duser.User, 0), nil
		}
		return nil, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	defer rows.Close()
	result := make([]*duser.User, 0)
	for rows.Next() {
		u, err := buildUserFromRow(rows, "")
		if err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, nil
}

func (r *PostgresUserRepository) Save(
	ctx context.Context,
	user *duser.User,
) error {
	exists, err := r.IDExists(ctx, user.ID())
	if err != nil {
		return err
	}
	if !exists {
		err = r.create(ctx, user)
		if err != nil {
			return err
		}
	}
	err = r.update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresUserRepository) create(
	ctx context.Context,
	user *duser.User,
) error {
	userStmt := fmt.Sprintf(
		`INSERT INTO %s (
			id,
			email,
			state,
			status,
			password_hash,
			version
		)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		userTable,
	)
	versionStmt := fmt.Sprintf(
		`INSERT INTO %s (
			user_id,
			email,
			state,
			status,
			password_hash,
			version
		)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		userVersionTable,
	)
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()
	_, err = tx.Exec(
		ctx,
		userStmt,
		user.ID(),
		user.Email().String(),
		user.State().String(),
		user.Status().String(),
		user.PasswordHash(),
		user.ModifyVersion(),
	)
	if err != nil {
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	_, err = tx.Exec(
		ctx,
		versionStmt,
		user.ID(),
		user.Email().String(),
		user.State().String(),
		user.Status().String(),
		user.PasswordHash(),
		user.ModifyVersion(),
	)
	if err != nil {
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	return nil
}

func (r *PostgresUserRepository) update(
	ctx context.Context,
	user *duser.User,
) error {
	query := fmt.Sprintf(
		`SELECT version
		FROM %s
		WHERE id = $1
		`,
		userTable,
	)
	userStmt := fmt.Sprintf(
		`UPDATE %s
		SET
			email = $2,
			state = $3,
			status = $4,
			password_hash = $5,
			version = $6
		WHERE id = $1`,
		userTable,
	)
	versionStmt := fmt.Sprintf(
		`INSERT INTO %s (
			user_id,
			email,
			state,
			status,
			password_hash,
			version
		)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		userVersionTable,
	)
	var currentVersion uint64
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()
	err = tx.QueryRow(ctx, query, user.ID()).Scan(&currentVersion)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf(
				"%w: пользователя с id %s не существует",
				application.ErrNotFound,
				user.ID(),
			)
		}
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	if currentVersion != user.Version() {
		return fmt.Errorf(
			"%w: версия пользователя с id %s не совпадает",
			application.ErrVersionConflict,
			user.ID(),
		)
	}
	_, err = tx.Exec(
		ctx,
		userStmt,
		user.ID(),
		user.Email().String(),
		user.State().String(),
		user.Status().String(),
		user.PasswordHash(),
		user.ModifyVersion(),
	)
	if err != nil {
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	_, err = tx.Exec(
		ctx,
		versionStmt,
		user.ID(),
		user.Email().String(),
		user.State().String(),
		user.Status().String(),
		user.PasswordHash(),
		user.ModifyVersion(),
	)
	if err != nil {
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	return nil
}

func buildUserFromRow(
	row pgx.Row, msg string,
) (*duser.User, error) {
	var (
		userID       uuid.UUID
		email        string
		state        string
		status       string
		passwordHash string
		version      uint64
	)
	if err := row.Scan(
		&userID,
		&email,
		&state,
		&status,
		&passwordHash,
		&version,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(
				"%w: %s",
				application.ErrNotFound,
				msg,
			)
		}
		return nil, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	domainEmail, err := domain.NewEmail(email)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	domainState, err := duser.StateFromString(state)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	domainStatus, err := duser.StatusFromString(status)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: %s",
			application.ErrInternal,
			err,
		)
	}
	u, err := duser.Restore(
		userID,
		domainEmail,
		domainState,
		domainStatus,
		passwordHash,
		version,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: %s",
			application.ErrInternal,
			err,
		)
	}
	return u, nil
}
