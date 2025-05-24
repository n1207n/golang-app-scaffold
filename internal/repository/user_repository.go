package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yourusername/yourprojectname/db/sqlc"       // Adjust to your actual module path
	"github.com/yourusername/yourprojectname/internal/util" // For password hashing
)

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	GetUserByID(ctx context.Context, id int64) (sqlc.User, error)
	GetUserByEmail(ctx context.Context, email string) (sqlc.User, error)
	// Add other methods like ListUsers, UpdateUser, DeleteUser as needed
}

// DBUserRepository implements UserRepository using a SQLC generated Querier.
type DBUserRepository struct {
	// Instead of db.Querier, sqlc generates a struct (default: Queries) that embeds this.
	// If you defined emit_interface: true and named it Store or Querier, use that.
	// For sqlc default, it's *db.Queries.
	// If you create a Store struct that embeds *db.Queries and also holds the *pgxpool.Pool,
	// then you'd pass the *pgxpool.Pool to NewStore and db.Store to NewDBUserRepository.
	// For simplicity, let's assume sqlc generates a `Store` interface or we use `*sqlc.Queries`.
	// Let's use `sqlc.Querier` as defined in sqlc.yaml (emit_interface: true)
	// and assume `db.New(dbpool)` returns something that implements `Querier`.
	// SQLC typically generates a `db.Queries` struct. If `emit_interface: true` is set,
	// it also generates a `Querier` interface.
	// A common pattern is to create a `Store` struct that embeds `*sqlc.Queries` and also holds the `*pgxpool.Pool`.
	// For now, let's assume `db.New(dbpool)` (from sqlc generated code) returns the querier.
	q *sqlc.Queries // This will be *sqlc.Queries if emit_interface is false or not customized
}

// NewDBUserRepository creates a new DBUserRepository.
// The `querier` argument should be the SQLC generated query executor, typically `*sqlc.Queries`.
func NewDBUserRepository(querier *sqlc.Queries) UserRepository {
	return &DBUserRepository{q: querier}
}

// CreateUser creates a new user in the database.
// It hashes the password before storing.
func (r *DBUserRepository) CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	hashedPassword, err := util.HashPassword(arg.HashedPassword) // Assume arg.HashedPassword is the plain password for now
	if err != nil {
		return sqlc.User{}, err
	}
	arg.HashedPassword = hashedPassword

	// Ensure Timestamptz fields are handled correctly if not defaulted by DB
	// sqlc.CreateUserParams might not have CreatedAt/UpdatedAt if DB defaults them.
	// If they are in params, they need to be set.
	// For this example, assuming DB defaults them or they are handled by sqlc based on query.

	return r.q.CreateUser(ctx, arg)
}

// GetUserByID retrieves a user by their ID.
func (r *DBUserRepository) GetUserByID(ctx context.Context, id int64) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

// GetUserByEmail retrieves a user by their email.
func (r *DBUserRepository) GetUserByEmail(ctx context.Context, email string) (sqlc.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

// TODO: Implement ListUsers, UpdateUser, DeleteUser
// Example for UpdateUser, note the use of pgtype for nullable fields if your sqlc query uses sqlc.narg()
func (r *DBUserRepository) UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
	// If password is being updated, it should be hashed
	if arg.HashedPassword.Valid { // Assuming HashedPassword is pgtype.Text for nullability
		hashedPassword, err := util.HashPassword(arg.HashedPassword.String)
		if err != nil {
			return sqlc.User{}, err
		}
		arg.HashedPassword = pgtype.Text{String: hashedPassword, Valid: true}
	}

	// Ensure UpdatedAt is set if your query expects it
	// arg.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true} // Example

	return r.q.UpdateUser(ctx, arg)
}
