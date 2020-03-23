package searchfield

type Workspace string

const (
	WorkspaceName     Workspace = "name"
	WorkspaceCreated  Workspace = "created"
	WorkspaceModified Workspace = "modified"
)

func (ws Workspace) String() string { return string(ws) }
