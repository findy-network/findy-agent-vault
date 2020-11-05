package model

import (
	"strconv"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type InternalPairwise struct {
	ID            string `faker:"uuid_hyphenated"`
	OurDid        string
	TheirDid      string
	TheirEndpoint string `faker:"url"`
	TheirLabel    string `faker:"organisationLabel"`
	InitiatedByUs bool
	ApprovedMs    int64 `faker:"unix_time"`
	CreatedMs     int64 `faker:"unix_time"`
}

func (p *InternalPairwise) Created() int64 {
	return p.CreatedMs
}

func (p *InternalPairwise) Identifier() string {
	return p.ID
}

func (p *InternalPairwise) Pairwise() *InternalPairwise {
	return p
}

func (p *InternalPairwise) Event() *InternalEvent {
	panic("Pairwise is not event")
}

func (p *InternalPairwise) Job() *InternalJob {
	panic("Pairwise is not job")
}

func (p *InternalPairwise) ToNode() *model.Pairwise {
	return &model.Pairwise{
		ID:            p.ID,
		OurDid:        p.OurDid,
		TheirDid:      p.TheirDid,
		TheirEndpoint: p.TheirEndpoint,
		TheirLabel:    p.TheirLabel,
		CreatedMs:     strconv.FormatInt(p.CreatedMs, 10),
		ApprovedMs:    strconv.FormatInt(p.ApprovedMs, 10),
		InitiatedByUs: p.InitiatedByUs,
	}
}

func (i *Items) PairwiseForID(id string) (node *model.Pairwise) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	if id == "" {
		return node
	}

	for _, item := range i.items {
		if item.Identifier() == id {
			node = item.Pairwise().ToNode()
			break
		}
	}

	return node
}

func (i *Items) PairwiseConnection(after, before int) *model.PairwiseConnection {
	i.mutex.RLock()
	result := i.items[after:before]
	totalCount := len(result)

	edges := make([]*model.PairwiseEdge, totalCount)
	nodes := make([]*model.Pairwise, totalCount)
	for index, pairwise := range result {
		node := pairwise.Pairwise().ToNode()
		edges[index] = &model.PairwiseEdge{
			Cursor: CreateCursor(pairwise.Pairwise().CreatedMs, model.Pairwise{}),
			Node:   node,
		}
		nodes[index] = node
	}
	i.mutex.RUnlock()

	var startCursor, endCursor *string
	if totalCount > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[totalCount-1].Cursor
	}
	p := &model.PairwiseConnection{
		Edges: edges,
		Nodes: nodes,
		PageInfo: &model.PageInfo{
			EndCursor:       endCursor,
			HasNextPage:     edges[len(edges)-1].Node.ID != i.LastID(),
			HasPreviousPage: edges[0].Node.ID != i.FirstID(),
			StartCursor:     startCursor,
		},
		TotalCount: totalCount,
	}
	return p
}
