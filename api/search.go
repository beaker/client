package api

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

type DatasetField string

const (
	DatasetAuthor            DatasetField = "user"
	DatasetCommitted         DatasetField = "committed"
	DatasetCreated           DatasetField = "created"
	DatasetName              DatasetField = "name"
	DatasetNameOrDescription DatasetField = "nameOrDescription"
)

func (ds DatasetField) String() string { return string(ds) }

type DatasetSearchOptions struct {
	SortClauses        []DatasetSortClause   `json:"sortClauses,omitempty"`
	FilterClauses      []DatasetFilterClause `json:"filterClauses,omitempty"`
	OmitResultDatasets bool                  `json:"omitResultDatasets,omitempty"`
	IncludeUncommitted bool                  `json:"includeUncommitted,omitempty"`
}

type DatasetSortClause struct {
	Field DatasetField `json:"field"`
	Order SortOrder    `json:"order"`
}

type DatasetFilterClause struct {
	Field    DatasetField   `json:"field"`
	Operator SearchOperator `json:"operator,omitempty"`
	Value    interface{}    `json:"value"`
}

type ExecutionField string

const (
	ExecutionID       ExecutionField = "id"
	ExecutionPriority ExecutionField = "priority"
)

func (e ExecutionField) String() string { return string(e) }

type ExperimentField string

const (
	ExperimentAuthor            ExperimentField = "author"
	ExperimentCreated           ExperimentField = "created"
	ExperimentName              ExperimentField = "name"
	ExperimentNameOrDescription ExperimentField = "nameOrDescription"
)

func (e ExperimentField) String() string { return string(e) }

type ExperimentSearchOptions struct {
	SortClauses   []ExperimentSortClause   `json:"sortClauses,omitempty"`
	FilterClauses []ExperimentFilterClause `json:"filterClauses,omitempty"`
}

type ExperimentSortClause struct {
	Field ExperimentField `json:"field"`
	Order SortOrder       `json:"order"`
}

type ExperimentFilterClause struct {
	Field    ExperimentField `json:"field"`
	Operator SearchOperator  `json:"operator,omitempty"`
	Value    interface{}     `json:"value"`
}

type GroupField string

const (
	GroupAuthor            GroupField = "author"
	GroupCreated           GroupField = "created"
	GroupModified          GroupField = "modified"
	GroupName              GroupField = "name"
	GroupNameOrDescription GroupField = "nameOrDescription"
)

func (g GroupField) String() string { return string(g) }

type GroupSearchOptions struct {
	SortClauses   []GroupSortClause   `json:"sortClauses,omitempty"`
	FilterClauses []GroupFilterClause `json:"filterClauses,omitempty"`
}

type GroupSortClause struct {
	Field GroupField `json:"field"`
	Order SortOrder  `json:"order"`
}

type GroupFilterClause struct {
	Field    GroupField     `json:"field"`
	Operator SearchOperator `json:"operator,omitempty"`
	Value    interface{}    `json:"value"`
}

type GroupTaskField string

const (
	GroupExperimentID GroupTaskField = "experimentId"
	GroupTaskID       GroupTaskField = "taskId"
)

func (gt GroupTaskField) String() string { return string(gt) }

type GroupTaskSearchOptions struct {
	SortClauses          []GroupTaskSortClause      `json:"sortClauses,omitempty"`
	ParameterSortClauses []GroupParameterSortClause `json:"parameterSortClauses,omitempty"`
}

type GroupTaskSortClause struct {
	Field GroupTaskField `json:"field"`
	Order SortOrder      `json:"order"`
}

type GroupParameterSortClause struct {
	Type  GroupParameterType `json:"type"`
	Name  string             `json:"name"`
	Order SortOrder          `json:"order"`
}

type ImageField string

const (
	ImageAuthor    ImageField = "author"
	ImageCommitted ImageField = "committed"
	ImageCreated   ImageField = "created"
	ImageName      ImageField = "name"
)

func (i ImageField) String() string { return string(i) }

type ImageSearchOptions struct {
	SortClauses   []ImageSortClause   `json:"sortClauses,omitempty"`
	FilterClauses []ImageFilterClause `json:"filterClauses,omitempty"`
}

type ImageSortClause struct {
	Field ImageField `json:"field"`
	Order SortOrder  `json:"order"`
}

type ImageFilterClause struct {
	Field    ImageField     `json:"field"`
	Operator SearchOperator `json:"operator,omitempty"`
	Value    interface{}    `json:"value"`
}

type WorkspaceField string

const (
	WorkspaceName     WorkspaceField = "name"
	WorkspaceCreated  WorkspaceField = "created"
	WorkspaceModified WorkspaceField = "modified"
)

func (ws WorkspaceField) String() string { return string(ws) }
