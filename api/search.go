package api

import (
	"github.com/beaker/client/api/searchfield"
)

type SearchOperator string

const (
	OpEqual            SearchOperator = "eq"
	OpGreaterThanEqual SearchOperator = "gte"
	OpLessThan         SearchOperator = "lt"
	OpContains         SearchOperator = "ctn"
)

type SortOrder string

const (
	SortAscending  SortOrder = "ascending"
	SortDescending SortOrder = "descending"
)

type ImageSearchOptions struct {
	SortClauses   []ImageSortClause   `json:"sortClauses,omitempty"`
	FilterClauses []ImageFilterClause `json:"filterClauses,omitempty"`
}

type ImageSortClause struct {
	Field searchfield.Image `json:"field"`
	Order SortOrder         `json:"order"`
}

type ImageFilterClause struct {
	Field    searchfield.Image `json:"field"`
	Operator SearchOperator    `json:"operator,omitempty"`
	Value    interface{}       `json:"value"`
}

type DatasetSearchOptions struct {
	SortClauses        []DatasetSortClause   `json:"sortClauses,omitempty"`
	FilterClauses      []DatasetFilterClause `json:"filterClauses,omitempty"`
	OmitResultDatasets bool                  `json:"omitResultDatasets,omitempty"`
	IncludeUncommitted bool                  `json:"includeUncommitted,omitempty"`
	Archived           *bool                 `json:"archived,omitempty"`
}

type DatasetSortClause struct {
	Field searchfield.Dataset `json:"field"`
	Order SortOrder           `json:"order"`
}

type DatasetFilterClause struct {
	Field    searchfield.Dataset `json:"field"`
	Operator SearchOperator      `json:"operator,omitempty"`
	Value    interface{}         `json:"value"`
}

type ExperimentSearchOptions struct {
	SortClauses   []ExperimentSortClause   `json:"sortClauses,omitempty"`
	FilterClauses []ExperimentFilterClause `json:"filterClauses,omitempty"`
	Archived      *bool                    `json:"archived,omitempty"`
}

type ExperimentSortClause struct {
	Field searchfield.Experiment `json:"field"`
	Order SortOrder              `json:"order"`
}

type ExperimentFilterClause struct {
	Field    searchfield.Experiment `json:"field"`
	Operator SearchOperator         `json:"operator,omitempty"`
	Value    interface{}            `json:"value"`
}

type GroupSearchOptions struct {
	SortClauses   []GroupSortClause   `json:"sortClauses,omitempty"`
	FilterClauses []GroupFilterClause `json:"filterClauses,omitempty"`
	Archived      *bool               `json:"archived,omitempty"`
}

type GroupSortClause struct {
	Field searchfield.Group `json:"field"`
	Order SortOrder         `json:"order"`
}

type GroupFilterClause struct {
	Field    searchfield.Group `json:"field"`
	Operator SearchOperator    `json:"operator,omitempty"`
	Value    interface{}       `json:"value"`
}

type GroupTaskSearchOptions struct {
	SortClauses          []GroupTaskSortClause      `json:"sortClauses,omitempty"`
	ParameterSortClauses []GroupParameterSortClause `json:"parameterSortClauses,omitempty"`
	FilterClauses        []GroupTaskFilterClause    `json:"filterClauses,omitempty"`
}

type GroupTaskSortClause struct {
	Field searchfield.GroupTask `json:"field"`
	Order SortOrder             `json:"order"`
}

type GroupParameterSortClause struct {
	Type  string    `json:"type"`
	Name  string    `json:"name"`
	Order SortOrder `json:"order"`
}

type GroupTaskFilterClause struct {
	Field    searchfield.GroupTask `json:"field"`
	Operator SearchOperator        `json:"operator,omitempty"`
	Value    interface{}           `json:"value"`
}
