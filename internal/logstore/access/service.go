package access

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/query"

	"github.com/zitadel/zitadel/internal/repository/quota"

	"github.com/zitadel/zitadel/internal/logstore"
)

type Config struct {
	Enabled  bool
	Debounce *logstore.DebouncerConfig
}

type Service struct {
	enabled   bool
	ctx       context.Context
	debouncer *logstore.Debouncer
	dbClient  *sql.DB
	report    reportFunc
}

type reportFunc func(ctx context.Context, q *query.Quota, used uint64) (doLimit bool, err error)

func NewService(ctx context.Context, cfg *Config, dbClient *sql.DB, report reportFunc) *Service {
	svc := &Service{
		enabled:  cfg != nil && cfg.Enabled,
		ctx:      ctx,
		dbClient: dbClient,
		report:   report,
	}

	if svc.enabled {
		if cfg.Debounce != nil && cfg.Debounce.MinFrequency > 0 && cfg.Debounce.MaxBulkSize > 0 {
			svc.debouncer = logstore.NewDebouncer(ctx, cfg.Debounce, newStorageBulkSink(dbClient))
		}
	}
	return svc
}

// TODO: Cache things in-memory here?
func (s *Service) Limit(ctx context.Context, instanceID string) (bool, error) {

	if instanceID == "" {
		return false, nil
	}

	quota, err := query.GetInstanceQuota(ctx, s.dbClient, instanceID, quota.RequestsAllAuthenticated)
	if err != nil {
		return false, err
	}

	usage, err := authenticatedInstanceRequests(ctx, s.dbClient, instanceID, quota.PeriodStart, quota.PeriodEnd)
	if err != nil {
		return false, err
	}

	return s.report(ctx, quota, usage)
}

func (s *Service) Handle(ctx context.Context, record *logstore.AccessLogRecord) (err error) {
	if !s.enabled {
		return nil
	}

	if s.debouncer != nil {
		s.debouncer.Add(record)
		return nil
	}

	return storeAccessLogs(ctx, s.dbClient, []any{record})
}

func (s *Service) Enabled() bool {
	return s.enabled
}