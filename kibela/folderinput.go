package kibela

// FolderInput represents the input folder with group of Kibela
type FolderInput struct {
	FolderName string `json:"folderName"`
	GroupId    ID     `json:"groupId"`
}
