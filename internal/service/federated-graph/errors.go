package federated_graph

const (
	ErrInvalidGraphName         = "Invalid Federated Graph Name"
	ErrCreatingGraph            = "Creating Federated Graph"
	ErrCompositionError         = "Composition Error"
	ErrRetrievingGraph          = "Error Retrieving Federated Graph"
	ErrInvalidResourceID        = "Invalid Resource ID"
	ErrReadingGraph             = "Error Reading Federated Graph"
	ErrUpdatingGraph            = "Error Updating Federated Graph"
	ErrDeletingGraph            = "Error Deleting Federated Graph"
	ErrUnexpectedDataSourceType = "Unexpected Data Source Configure Type"
	ErrCompositionErrors        = "Composition Errors"
)

const (
	DebugCreate = "create-federated-graph"
	DebugRead   = "read-federated-graph"
	DebugUpdate = "update-federated-graph"
	DebugDelete = "delete-federated-graph"
)
