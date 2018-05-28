package utils

import (
	"github.com/graphql-go/graphql"
	"github.com/jenarvaezg/mongodm"
)

// MakeListField returns a GraphQL Field that is a list of the listType element with limit and offset args
// if searchable, an extra argument, "search", is added
func MakeListField(listType graphql.Output, resolve graphql.FieldResolveFn, searchable bool) *graphql.Field {
	field := &graphql.Field{Type: listType, Resolve: resolve}
	args := graphql.FieldConfigArgument{
		"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
		"offset": &graphql.ArgumentConfig{Type: graphql.Int},
	}

	if searchable {
		args["search"] = &graphql.ArgumentConfig{Type: graphql.String}
	}

	field.Args = args
	return field
}

// MakeNodeListType returns a GraphQL object for a list of objects, with nodes and total count
func MakeNodeListType(name string, nodeType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: name,
		Fields: graphql.Fields{
			"nodes":      &graphql.Field{Type: graphql.NewList(nodeType)},
			"totalCount": &graphql.Field{Type: graphql.Int},
		},
	})
}

// MergeGraphQLFields receives a list of graphql.Fields objects and merges their values
func MergeGraphQLFields(grapqhQLFields ...graphql.Fields) graphql.Fields {
	mergedFields := graphql.Fields{}
	for _, f := range grapqhQLFields {
		for k, v := range f {
			mergedFields[k] = v
		}
	}

	return mergedFields
}

// GetIDResolver resolves the ID of a document
func GetIDResolver(p graphql.ResolveParams) (interface{}, error) {
	team := p.Source.(mongodm.IDocumentBase)
	return team.GetId().Hex(), nil
}
