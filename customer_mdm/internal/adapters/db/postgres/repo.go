package postgres

import (
    "context"
    "errors"
    "strings"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/yourorg/customer-mdm/internal/domain"
    "github.com/yourorg/customer-mdm/internal/ports"
    "github.com/google/uuid"
)

type CustomerRepository struct {
    Pool *pgxpool.Pool
}

func (r *CustomerRepository) Create(ctx context.Context, c *domain.Customer) error {
    if c.ID == "" {
        c.ID = uuid.New().String()
    }
    now := time.Now().UTC()
    c.CreatedAt, c.UpdatedAt = now, now
    const q = `INSERT INTO customers (id,name,email,tax_id,phone,country,is_active,attributes,created_at,updated_at)
               VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
    _, err := r.Pool.Exec(ctx, q, c.ID, c.Name, c.Email, c.TaxID, c.Phone, c.Country, c.IsActive, c.Attributes, c.CreatedAt, c.UpdatedAt)
    return err
}

func (r *CustomerRepository) Get(ctx context.Context, id string) (*domain.Customer, error) {
    const q = `SELECT id,name,email,tax_id,phone,country,is_active,attributes,created_at,updated_at FROM customers WHERE id=$1`
    row := r.Pool.QueryRow(ctx, q, id)
    var c domain.Customer
    if err := row.Scan(&c.ID,&c.Name,&c.Email,&c.TaxID,&c.Phone,&c.Country,&c.IsActive,&c.Attributes,&c.CreatedAt,&c.UpdatedAt); err != nil {
        if errors.Is(err, pgx.ErrNoRows) { return nil, nil }
        return nil, err
    }
    return &c, nil
}

func (r *CustomerRepository) Update(ctx context.Context, c *domain.Customer) error {
    c.UpdatedAt = time.Now().UTC()
    const q = `UPDATE customers SET name=$2,email=$3,tax_id=$4,phone=$5,country=$6,is_active=$7,attributes=$8,updated_at=$9 WHERE id=$1`
    cmd, err := r.Pool.Exec(ctx, q, c.ID, c.Name, c.Email, c.TaxID, c.Phone, c.Country, c.IsActive, c.Attributes, c.UpdatedAt)
    if err != nil { return err }
    if cmd.RowsAffected() == 0 { return pgx.ErrNoRows }
    return nil
}

func (r *CustomerRepository) Delete(ctx context.Context, id string) error {
    const q = `DELETE FROM customers WHERE id=$1`
    _, err := r.Pool.Exec(ctx, q, id)
    return err
}

func (r *CustomerRepository) List(ctx context.Context, pageSize int, pageToken, query string) (items []*domain.Customer, next string, total int, err error) {
    if pageSize <= 0 { pageSize = 50 }
    if pageSize > 500 { pageSize = 500 }

    where := ""
    args := []any{pageSize}
    if q := strings.TrimSpace(query); q != "" {
        where = "WHERE name ILIKE '%' || $2 || '%' OR email ILIKE '%' || $2 || '%'"
        args = append([]any{pageSize, q}, args[1:]...)
    }
    sql := `SELECT id,name,email,tax_id,phone,country,is_active,attributes,created_at,updated_at
            FROM customers ` + where + `
            ORDER BY created_at DESC
            LIMIT $1`
    rows, err := r.Pool.Query(ctx, sql, args...)
    if err != nil { return nil, "", 0, err }
    defer rows.Close()
    for rows.Next() {
        var c domain.Customer
        if err := rows.Scan(&c.ID,&c.Name,&c.Email,&c.TaxID,&c.Phone,&c.Country,&c.IsActive,&c.Attributes,&c.CreatedAt,&c.UpdatedAt); err != nil {
            return nil, "", 0, err
        }
        items = append(items, &c)
    }
    // optional count
    _ = r.Pool.QueryRow(ctx, "SELECT COUNT(1) FROM customers").Scan(&total)
    next = "" // cursor omitted for brevity
    return
}

type OutboxWriter struct {
    Pool *pgxpool.Pool
}

func (o *OutboxWriter) Write(tx ports.Tx, topic string, key, value []byte, headers map[string]string) error {
    const q = `INSERT INTO outbox_events (aggregate_type, aggregate_id, event_type, payload, headers)
               VALUES ($1,$2,$3,$4,$5)`
    // we expect value is a protobuf CustomerEvent; we parse minimal to get aggregate_id/event_type? For simplicity pass placeholders.
    aggregateType := "customer"
    aggregateID := uuid.New().String()
    eventType := "generic"
    _, err := o.Pool.Exec(tx.Context(), q, aggregateType, aggregateID, eventType, value, headers)
    return err
}
