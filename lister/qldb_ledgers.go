package lister

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/qldb"
	"github.com/trek10inc/awsets/context"
	"github.com/trek10inc/awsets/resource"
)

type AWSQLDBLedger struct {
}

func init() {
	i := AWSQLDBLedger{}
	listers = append(listers, i)
}

func (l AWSQLDBLedger) Types() []resource.ResourceType {
	return []resource.ResourceType{resource.QLDBLedger}
}

func (l AWSQLDBLedger) List(ctx context.AWSetsCtx) (*resource.Group, error) {
	svc := qldb.New(ctx.AWSCfg)

	req := svc.ListLedgersRequest(&qldb.ListLedgersInput{
		MaxResults: aws.Int64(100),
	})

	rg := resource.NewGroup()
	paginator := qldb.NewListLedgersPaginator(req)
	for paginator.Next(ctx.Context) {
		page := paginator.CurrentPage()
		for _, v := range page.Ledgers {
			r := resource.New(ctx, resource.QLDBLedger, v.Name, v.Name, v)
			rg.AddResource(r)
		}
	}
	err := paginator.Err()
	return rg, err
}
