package repository

import (
	"context"
	"erp_monolith/backend/internal/modules/system/domain"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresCompanyRepo struct {
	pool *pgxpool.Pool
}

// NewPostgresCompanyRepo creates a driven adapter repository
func NewPostgresCompanyRepo(pool *pgxpool.Pool) domain.CompanyRepository {
	return &postgresCompanyRepo{pool: pool}
}

func (r *postgresCompanyRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	query := `SELECT id, code, name, currency, is_active, created_at, updated_at FROM adm_companies WHERE id = $1`
	var c domain.Company
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Code, &c.Name, &c.Currency, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query company: %w", err)
	}
	return &c, nil
}

func (r *postgresCompanyRepo) FindAll(ctx context.Context) ([]domain.Company, error) {
	query := `SELECT id, code, name, currency, is_active, created_at, updated_at FROM adm_companies ORDER BY code ASC`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list companies: %w", err)
	}
	defer rows.Close()

	var companies []domain.Company
	for rows.Next() {
		var c domain.Company
		if err := rows.Scan(&c.ID, &c.Code, &c.Name, &c.Currency, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, nil
}

func (r *postgresCompanyRepo) Save(ctx context.Context, company *domain.Company) error {
	query := `
		INSERT INTO adm_companies (id, code, name, currency, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			code = EXCLUDED.code,
			name = EXCLUDED.name,
			currency = EXCLUDED.currency,
			is_active = EXCLUDED.is_active,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.pool.Exec(ctx, query,
		company.ID, company.Code, company.Name, company.Currency, company.IsActive, company.CreatedAt, company.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save company: %w", err)
	}
	return nil
}
