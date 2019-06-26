package api

import (
	"github.com/beaker/client/api/searchfield"
)

type SearchOperator string

const (
	OpEqual            SearchOperator = "eq"
	OpNotEqual         SearchOperator = "neq"
	OpGreaterThan      SearchOperator = "gt"
	OpGreaterThanEqual SearchOperator = "gte"
	OpLessThan         SearchOperator = "lt"
	OpLessThanEqual    SearchOperator = "lte"
	OpContains         SearchOperator = "ctn"
	OpNotContains      SearchOperator = "nctn"
)

type SortOrder string

const (
	SortAscending  SortOrder = "ascending"
	SortDescending SortOrder = "descending"
)

type FilterCombinator string

const (
	CombinatorAnd FilterCombinator = "and"
	CombinatorOr  FilterCombinator = "or"
)

type ImageSearchOptions struct {
	SortClauses                []ImageSortClause   `json:"sortClauses,omitempty"`
	SortClausesDeprecated      []ImageSortClause   `json:"sort_clauses,omitempty"`
	FilterClauses              []ImageFilterClause `json:"filterClauses,omitempty"`
	FilterClausesDeprecated    []ImageFilterClause `json:"filter_clauses,omitempty"`
	FilterCombinator           FilterCombinator    `json:"filterCombinator,omitempty"`
	FilterCombinatorDeprecated FilterCombinator    `json:"filter_combinator,omitempty"`
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
	SortClauses                  []DatasetSortClause   `json:"sortClauses,omitempty"`
	SortClausesDeprecated        []DatasetSortClause   `json:"sort_clauses,omitempty"`
	FilterClauses                []DatasetFilterClause `json:"filterClauses,omitempty"`
	FilterClausesDeprecated      []DatasetFilterClause `json:"filter_clauses,omitempty"`
	FilterCombinator             FilterCombinator      `json:"filterCombinator,omitempty"`
	FilterCombinatorDeprecated   FilterCombinator      `json:"filter_combinator,omitempty"`
	OmitResultDatasets           bool                  `json:"omitResultDatasets,omitempty"`
	OmitResultDatasetsDeprecated bool                  `json:"omit_result_datasets,omitempty"`
	IncludeUncommitted           bool                  `json:"includeUncommitted,omitempty"`
	IncludeUncommittedDeprecated bool                  `json:"include_uncommitted,omitempty"`
	Archived                     *bool                 `json:"archived,omitempty"`
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
	SortClauses                []ExperimentSortClause   `json:"sortClauses,omitempty"`
	SortClausesDeprecated      []ExperimentSortClause   `json:"sort_clauses,omitempty"`
	FilterClauses              []ExperimentFilterClause `json:"filterClauses,omitempty"`
	FilterClausesDeprecated    []ExperimentFilterClause `json:"filter_clauses,omitempty"`
	FilterCombinator           FilterCombinator         `json:"filterCombinator,omitempty"`
	FilterCombinatorDeprecated FilterCombinator         `json:"filter_combinator,omitempty"`
	Archived                   *bool                    `json:"archived,omitempty"`
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
	SortClauses                []GroupSortClause   `json:"sortClauses,omitempty"`
	SortClausesDeprecated      []GroupSortClause   `json:"sort_clauses,omitempty"`
	FilterClauses              []GroupFilterClause `json:"filterClauses,omitempty"`
	FilterClausesDeprecated    []GroupFilterClause `json:"filter_clauses,omitempty"`
	FilterCombinator           FilterCombinator    `json:"filterCombinator,omitempty"`
	FilterCombinatorDeprecated FilterCombinator    `json:"filter_combinator,omitempty"`
	Archived                   *bool               `json:"archived,omitempty"`
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
	SortClauses                    []GroupTaskSortClause      `json:"sortClauses,omitempty"`
	SortClausesDeprecated          []GroupTaskSortClause      `json:"sort_clauses,omitempty"`
	ParameterSortClauses           []GroupParameterSortClause `json:"parameterSortClauses,omitempty"`
	ParameterSortClausesDeprecated []GroupParameterSortClause `json:"parameter_sort_clauses,omitempty"`
	FilterClauses                  []GroupTaskFilterClause    `json:"filterClauses,omitempty"`
	FilterClausesDeprecated        []GroupTaskFilterClause    `json:"filter_clauses,omitempty"`
	FilterCombinator               FilterCombinator           `json:"filterCombinator,omitempty"`
	FilterCombinatorDeprecated     FilterCombinator           `json:"filter_combinator,omitempty"`
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
