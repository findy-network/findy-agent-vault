package resolver

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-vault/graph/model"
	data "github.com/findy-network/findy-agent-vault/tools/data/model"
	"github.com/golang/glog"
	"github.com/lainio/err2"
)

func (r *credentialResolver) Connection(ctx context.Context, obj *model.Credential) (c *model.Pairwise, err error) {
	glog.V(logLevelMedium).Info("credentialResolver:Connection, id: ", obj.ID)
	defer err2.Return(&err)

	if connectionID := state.Credentials().CredentialPairwiseID(obj.ID); connectionID != nil {
		return r.Query().Connection(ctx, *connectionID)
	}

	err = fmt.Errorf("pairwise for credential id %s was not found", obj.ID)
	return
}

func (r *queryResolver) Credential(ctx context.Context, id string) (node *model.Credential, err error) {
	glog.V(logLevelMedium).Info("queryResolver:Credential, id: ", id)

	items := state.Credentials()
	edge := items.CredentialForID(id)
	if edge == nil {
		err = fmt.Errorf("connection for id %s was not found", id)
	} else {
		node = edge.Node
	}
	return
}

func (r *queryResolver) Credentials(
	ctx context.Context,
	after, before *string,
	first, last *int,
) (c *model.CredentialConnection, err error) {
	defer err2.Return(&err)

	pagination := &PaginationParams{
		first:  first,
		last:   last,
		after:  after,
		before: before,
	}
	logPaginationRequest("queryResolver:credentials", pagination)

	items := state.Credentials()
	items = &data.CredentialItems{Items: items.Filter(func(item data.APIObject) data.APIObject {
		c := item.Credential()
		fmt.Println(c)
		if c.IssuedMs != nil {
			return c.Copy()
		}
		return nil
	}),
	}

	afterIndex, beforeIndex, err := pick(items.Objects(), pagination)
	err2.Check(err)

	glog.V(logLevelLow).Infof("Credentials: returning connections between %d and %d", afterIndex, beforeIndex)
	c = items.CredentialConnection(afterIndex, beforeIndex)

	return c, err
}

func (r *pairwiseResolver) Credentials(
	ctx context.Context,
	obj *model.Pairwise,
	after, before *string,
	first, last *int,
) (c *model.CredentialConnection, err error) {
	defer err2.Return(&err)
	pagination := &PaginationParams{
		first:  first,
		last:   last,
		after:  after,
		before: before,
	}
	logPaginationRequest("pairwiseResolver:credentials", pagination)

	items := state.Credentials()
	items = &data.CredentialItems{Items: items.Filter(func(item data.APIObject) data.APIObject {
		c := item.Credential()
		fmt.Println(c)
		if c.IssuedMs != nil && c.PairwiseID == obj.ID {
			return c.Copy()
		}
		return nil
	}),
	}

	afterIndex, beforeIndex, err := pick(items.Objects(), pagination)
	err2.Check(err)

	glog.V(logLevelLow).Infof("Credentials: returning credentials between %d and %d", afterIndex, beforeIndex)

	return items.CredentialConnection(afterIndex, beforeIndex), nil
}
