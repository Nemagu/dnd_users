package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	userTable        = "\"user\""
	userVersionTable = "user_version"
)

type PostgresUserRepository struct {
	logger *slog.Logger
	pool   *pgxpool.Pool
}

func NewPostgresUserRepository(
	logger *slog.Logger,
	pool *pgxpool.Pool,
) (*PostgresUserRepository, error) {
	return &PostgresUserRepository{
		logger: logger,
		pool:   pool,
	}, nil
}

func MustNewPostgresUserRepository(
	logger *slog.Logger,
	pool *pgxpool.Pool,
) *PostgresUserRepository {
	return &PostgresUserRepository{
		logger: logger,
		pool:   pool,
	}
}

func (r *PostgresUserRepository) NextID(ctx context.Context) uuid.UUID {
	return uuid.New()
}

func (r *PostgresUserRepository) IDExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	query := fmt.Sprintf(
		"SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)",
		userTable,
	)
	r.logger.DebugContext(ctx, "create query", "query", query)
	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return exists, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	r.logger.DebugContext(ctx, "executed query", "query", query, "result", exists)
	return exists, nil
}

func (r *PostgresUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := fmt.Sprintf(
		"SELECT EXISTS(SELECT 1 FROM %s WHERE email = $1)",
		userTable,
	)
	r.logger.InfoContext(ctx, "try to execute query", "query", query)
	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return exists, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	r.logger.DebugContext(ctx, "create query", "query", query, "result", exists)
	return exists, nil
}

func (r *PostgresUserRepository) All(
	ctx context.Context,
	limit, offset int,
) ([]*appdto.User, error) {
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
		r.logger.DebugContext(ctx, "add limit to query", "limit", limit)
		query += fmt.Sprintf(" LIMIT $%d", len(values)+1)
		values = append(values, limit)
	}
	if offset > 0 {
		r.logger.DebugContext(ctx, "add offset to query", "offset", offset)
		query += fmt.Sprintf(" OFFSET $%d", len(values)+1)
		values = append(values, offset)
	}
	query += " ORDER BY email"
	r.logger.InfoContext(ctx, "try to execute query", "query", query)
	rows, err := r.pool.Query(ctx, query, values...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return make([]*appdto.User, 0), nil
		}
		return nil, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	defer rows.Close()
	cap := 30
	if limit != 0 {
		cap = limit
	}
	result := make([]*appdto.User, 0, cap)
	r.logger.InfoContext(ctx, "generate dto from executed query")
	for rows.Next() {
		u, err := buildUserFromRow(rows, "")
		if err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	r.logger.DebugContext(ctx, "success generated dto", "dto", result)
	return result, nil
}

func (r *PostgresUserRepository) ByID(ctx context.Context, id uuid.UUID) (*appdto.User, error) {
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
	r.logger.InfoContext(ctx, "try to execute query", "query", query)
	row := r.pool.QueryRow(ctx, query, id)
	r.logger.InfoContext(ctx, "try generate dto from executed dto")
	return buildUserFromRow(
		row,
		fmt.Sprintf("пользователя с id %s не существует", id),
	)
}

func (r *PostgresUserRepository) ByEmail(ctx context.Context, email string) (*appdto.User, error) {
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
	r.logger.InfoContext(ctx, "try to execute query", "query", query)
	row := r.pool.QueryRow(ctx, query, email)
	r.logger.InfoContext(ctx, "try generate dto from executed dto")
	return buildUserFromRow(
		row,
		fmt.Sprintf("пользователя с email %s не существует", email),
	)
}

func (r *PostgresUserRepository) Filter(
	ctx context.Context,
	searchByEmail string,
	filterByState []string,
	filterByStatus []string,
	limit, offset int,
) ([]*appdto.User, error) {
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
		r.logger.DebugContext(ctx, "add searching by email to query", "email", searchByEmail)
		condition := fmt.Sprintf("LOWER(email) LIKE $%d", len(values)+1)
		conditions = append(conditions, condition)
		values = append(values, "%"+strings.ToLower(searchByEmail)+"%")
	}
	if len(filterByState) != 0 {
		r.logger.DebugContext(ctx, "add filtering by state to query", "states", filterByState)
		condition := "state IN (" + strings.Join(filterByState, ",") + ")"
		conditions = append(conditions, condition)
	}
	if len(filterByStatus) != 0 {
		r.logger.DebugContext(ctx, "add filtering by status to query", "statuses", filterByStatus)
		condition := "status IN (" + strings.Join(filterByStatus, ",") + ")"
		conditions = append(conditions, condition)
	}
	if len(values) > 0 {
		r.logger.DebugContext(ctx, "join conditions")
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	if limit > 0 {
		r.logger.DebugContext(ctx, "add limit to query", "limit", limit)
		query += fmt.Sprintf(" LIMIT $%d", len(values)+1)
		values = append(values, limit)
	}
	if offset > 0 {
		r.logger.DebugContext(ctx, "add offset to query", "offset", offset)
		query += fmt.Sprintf(" OFFSET $%d", len(values)+1)
		values = append(values, offset)
	}
	query += " ORDER BY email"
	r.logger.InfoContext(ctx, "try to execute query", "query", query)
	rows, err := r.pool.Query(ctx, query, values...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return make([]*appdto.User, 0), nil
		}
		return nil, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	defer rows.Close()
	result := make([]*appdto.User, 0)
	r.logger.InfoContext(ctx, "try generate dto from executed dto")
	for rows.Next() {
		u, err := buildUserFromRow(rows, "")
		if err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	r.logger.DebugContext(ctx, "success generated dto", "dto", result)
	return result, nil
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *appdto.User) error {
	r.logger.DebugContext(ctx, "check exists")
	exists, err := r.IDExists(ctx, user.UserID)
	if err != nil {
		return err
	}
	switch {
	case exists:
		r.logger.InfoContext(ctx, "try to create user")
		err = r.create(ctx, user)
		if err != nil {
			return err
		}
		r.logger.DebugContext(ctx, "user created", "user", user)
	default:
		r.logger.InfoContext(ctx, "try to update user")
		err = r.update(ctx, user)
		if err != nil {
			return err
		}
		r.logger.DebugContext(ctx, "user updated", "user", user)
	}
	return nil
}

func (r *PostgresUserRepository) create(ctx context.Context, user *appdto.User) error {
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
	r.logger.InfoContext(ctx, "try to begin transaction")
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	defer func() {
		if tx != nil {
			if err := tx.Rollback(ctx); err != nil {
				r.logger.ErrorContext(ctx, "rollback failed")
			}
		}
	}()
	r.logger.InfoContext(ctx, "try to execute statement", "stmt", userStmt)
	_, err = tx.Exec(
		ctx,
		userStmt,
		user.UserID,
		user.Email,
		user.State,
		user.Status,
		user.PasswordHash,
		user.Version,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	r.logger.InfoContext(ctx, "try to execute statement", "stmt", versionStmt)
	_, err = tx.Exec(
		ctx,
		versionStmt,
		user.UserID,
		user.Email,
		user.State,
		user.Status,
		user.PasswordHash,
		user.Version,
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

func (r *PostgresUserRepository) update(ctx context.Context, user *appdto.User) error {
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
	r.logger.InfoContext(ctx, "try to begin transaction")
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	defer func() {
		if tx != nil {
			if err := tx.Rollback(ctx); err != nil {
				r.logger.ErrorContext(ctx, "rollback failed")
			}
		}
	}()
	r.logger.InfoContext(ctx, "try to execute query", "query", query)
	err = tx.QueryRow(ctx, query, user.UserID).Scan(&currentVersion)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf(
				"%w: пользователя с id %s не существует",
				application.ErrNotFound,
				user.UserID,
			)
		}
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	if currentVersion >= user.Version {
		return fmt.Errorf(
			"%w: версия пользователя с id %s не совпадает",
			application.ErrVersionConflict,
			user.UserID,
		)
	}
	r.logger.InfoContext(ctx, "try to execute statement", "stmt", userStmt)
	_, err = tx.Exec(
		ctx,
		userStmt,
		user.UserID,
		user.Email,
		user.State,
		user.Status,
		user.PasswordHash,
		user.Version,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	r.logger.InfoContext(ctx, "try to execute statement", "stmt", versionStmt)
	_, err = tx.Exec(
		ctx,
		versionStmt,
		user.UserID,
		user.Email,
		user.State,
		user.Status,
		user.PasswordHash,
		user.Version,
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

func buildUserFromRow(row pgx.Row, msg string) (*appdto.User, error) {
	u := appdto.User{}
	if err := row.Scan(
		&u.UserID,
		&u.Email,
		&u.State,
		&u.Status,
		&u.PasswordHash,
		&u.Version,
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
	return &u, nil
}
