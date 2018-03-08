package migrations

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

var GlobalMigrations *Migrations

func initGM() {
	if GlobalMigrations == nil {
		GlobalMigrations = &Migrations{migrations: make(map[string]Migration)}
	}
}

type Migration interface {
	Sequence() string
	Migrate(ctx context.Context, execer sqlx.Execer) context.Context
	IsApplied(ctx context.Context, q sqlx.Queryer) bool
}

type Migrations struct {
	migrations map[string]Migration
}

func (m *Migrations) Apply(tx *sqlx.DB) error {
	ctx := context.Background()
	ctx = EnsureMigrationsTable(ctx, tx)
	var migrationSeqs []string
	for k, _ := range m.migrations {
		migrationSeqs = append(migrationSeqs, k)
	}
	sort.Strings(migrationSeqs)
	for _, seq := range migrationSeqs {
		migration := m.migrations[seq]
		if !migration.IsApplied(ctx, tx) {
			ctx = m.migrations[seq].Migrate(ctx, tx)
			tx.MustExec(`insert into database_migrations (sequence, applied_at) values ($1, $2)`, migration.Sequence(), time.Now())
		}
	}
	if ctx.Value("error") != nil {
		log.Fatal("Error while running migration at stage ... ", ctx.Value("stage"))
		log.Fatal(ctx.Value("error"))
		return ctx.Value("error").(error)
	}
	return nil
}

func (m *Migrations) Register(migration Migration) {
	m.migrations[migration.Sequence()] = migration
}

const MigrationsSQL = `
	create table if not exists database_migrations(
		sequence text primary key,
		applied_at timestamp
	)
`

type DatabaseMigration struct {
	Sequence  string    `db:"sequence"`
	AppliedAt time.Time `db:"applied_at"`
}

func EnsureMigrationsTable(ctx context.Context, execer sqlx.Execer) context.Context {
	if ctx.Value("error") != nil {
		return ctx
	}
	_, err := execer.Exec(MigrationsSQL)
	if err != nil {
		log.Fatal(err)
		errCtx := context.WithValue(ctx, "error", err)
		return context.WithValue(errCtx, "stage", "EnsureMigrationsTable")
	}
	return ctx
}

type SQLMigration struct {
	sequence string
	sql      string
}

func NewMigration(sequence, sql string) *SQLMigration {
	return &SQLMigration{sequence: sequence, sql: sql}
}

func (m *SQLMigration) Sequence() string {
	return m.sequence
}
func (m *SQLMigration) IsApplied(ctx context.Context, q sqlx.Queryer) bool {
	if ctx.Value("error") != nil {
		return false
	}
	row := q.QueryRowx("select sequence, applied_at from database_migrations where sequence=$1", m.sequence)
	var dbm DatabaseMigration
	err := row.StructScan(&dbm)
	if err != nil {
		return false
	}
	return true
}

func (m *SQLMigration) Migrate(ctx context.Context, execer sqlx.Execer) context.Context {
	if ctx.Value("error") != nil {
		return ctx
	}
	_, err := execer.Exec(m.sql)
	if err != nil {
		log.Fatal(err)
		errCtx := context.WithValue(ctx, "error", err)
		return context.WithValue(errCtx, "stage", fmt.Sprintf("RunMigration - %s", m.sequence))
	}
	return ctx
}
