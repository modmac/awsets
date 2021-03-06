package lister

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/trek10inc/awsets/context"
	"github.com/trek10inc/awsets/resource"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type AWSNeptuneDbClusterSnapshot struct {
}

func init() {
	i := AWSNeptuneDbClusterSnapshot{}
	listers = append(listers, i)
}

func (l AWSNeptuneDbClusterSnapshot) Types() []resource.ResourceType {
	return []resource.ResourceType{resource.NeptuneDbClusterSnapshot}
}

func (l AWSNeptuneDbClusterSnapshot) List(ctx context.AWSetsCtx) (*resource.Group, error) {
	svc := neptune.New(ctx.AWSCfg)

	rg := resource.NewGroup()

	var marker *string

	for {
		res, err := svc.DescribeDBClusterSnapshotsRequest(&neptune.DescribeDBClusterSnapshotsInput{
			Marker:     marker,
			MaxRecords: aws.Int64(100),
		}).Send(ctx.Context)
		if err != nil {
			return rg, fmt.Errorf("failed to list neptune cluster snapshots: %w", err)
		}
		for _, v := range res.DBClusterSnapshots {
			r := resource.New(ctx, resource.NeptuneDbClusterSnapshot, v.DBClusterSnapshotIdentifier, v.DBClusterSnapshotIdentifier, v)
			r.AddARNRelation(resource.KmsKey, v.KmsKeyId)
			r.AddRelation(resource.NeptuneDbCluster, v.DBClusterIdentifier, "")
			r.AddRelation(resource.Ec2Vpc, v.VpcId, "")

			rg.AddResource(r)
		}
		if res.Marker == nil {
			break
		}
		marker = res.Marker
	}
	return rg, nil
}
