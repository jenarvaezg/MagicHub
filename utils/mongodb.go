package utils

import "github.com/zebresel-com/mongodm"

// QueryLimitAndOffset Receives limit, offset and a mongodb query, returns the mondodm query, ready to be executed
func QueryLimitAndOffset(limit, offset int, query *mongodm.Query) *mongodm.Query {
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Skip(offset)
	}
	return query
}
