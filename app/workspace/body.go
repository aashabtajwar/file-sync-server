package workspace

type WorkspaceCreationData struct {
	Name string
}

type WorkspaceFiles struct {
	Files []string
}

type UserEmail struct {
	Email string
}

type PayLoad struct {
	Message string
	Id      string
}
