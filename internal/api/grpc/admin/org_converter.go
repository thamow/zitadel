package admin

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	org_grpc "github.com/caos/zitadel/internal/api/grpc/org"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/pkg/grpc/admin"
)

func listOrgRequestToModel(req *admin.ListOrgsRequest) (*query.OrgSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := org_grpc.OrgQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.OrgSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			SortingColumn: req.SortingColumn.String(),
			Asc:           asc,
		},
		Queries: queries,
	}, nil
}

func setUpOrgOrgToDomain(req *admin.SetUpOrgRequest_Org) *domain.Org {
	org := &domain.Org{
		Name:    req.Name,
		Domains: []*domain.OrgDomain{},
	}
	if req.Domain != "" {
		org.Domains = append(org.Domains, &domain.OrgDomain{Domain: req.Domain})
	}
	return org
}
