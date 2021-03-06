package lister

import (
	"fmt"
	"sync"

	"github.com/trek10inc/awsets/context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/waf"
	"github.com/trek10inc/awsets/resource"
)

var listWafByteMatchSetsOnce sync.Once

type AWSWafByteMatchSet struct {
}

func init() {
	i := AWSWafByteMatchSet{}
	listers = append(listers, i)
}

func (l AWSWafByteMatchSet) Types() []resource.ResourceType {
	return []resource.ResourceType{resource.WafByteMatchSet}
}

func (l AWSWafByteMatchSet) List(ctx context.AWSetsCtx) (*resource.Group, error) {
	svc := waf.New(ctx.AWSCfg)
	rg := resource.NewGroup()

	var outerErr error

	listWafByteMatchSetsOnce.Do(func() {
		var nextMarker *string
		for {
			res, err := svc.ListByteMatchSetsRequest(&waf.ListByteMatchSetsInput{
				Limit:      aws.Int64(100),
				NextMarker: nextMarker,
			}).Send(ctx.Context)
			if err != nil {
				outerErr = fmt.Errorf("failed to list byte match sets: %w", err)
				return
			}
			for _, id := range res.ByteMatchSets {
				byteMatchSet, err := svc.GetByteMatchSetRequest(&waf.GetByteMatchSetInput{
					ByteMatchSetId: id.ByteMatchSetId,
				}).Send(ctx.Context)
				if err != nil {
					outerErr = fmt.Errorf("failed to get byte match stringset %s: %w", aws.StringValue(id.ByteMatchSetId), err)
					return
				}
				if v := byteMatchSet.ByteMatchSet; v != nil {
					r := resource.NewGlobal(ctx, resource.WafByteMatchSet, v.ByteMatchSetId, v.Name, v)
					rg.AddResource(r)
				}
			}
			if res.NextMarker == nil {
				break
			}
			nextMarker = res.NextMarker
		}
	})

	return rg, outerErr
}
