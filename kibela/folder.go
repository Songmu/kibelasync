package kibela

import (
	"encoding/json"

	"github.com/Songmu/kibelasync/client"
	"golang.org/x/xerrors"
)

// Folder represents folder of Kibela
type Folder struct {
	ID       `json:"id"`
	FullName string `json:"fullName"`
}

func (ki *Kibela) getFolderCount() (int, error) {
	data, err := ki.cli.Do(&client.Payload{Query: totalFolderCountQuery})
	if err != nil {
		return 0, xerrors.Errorf("failed to ki.getFolderCount: %w", err)
	}
	var res struct {
		Folders struct {
			TotalCount int `json:"totalCount"`
		} `json:"folders"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return 0, xerrors.Errorf("failed to ki.getNotesCount: %w", err)
	}
	return res.Folders.TotalCount, nil
}

func (ki *Kibela) getFolders() ([]*Folder, error) {
	num, err := ki.getFolderCount()
	if err != nil {
		return nil, xerrors.Errorf("failed to getFolders: %w", err)
	}
	data, err := ki.cli.Do(&client.Payload{Query: listFolderQuery(num)})
	if err != nil {
		return nil, xerrors.Errorf("failed to ki.getFolders: %w", err)
	}
	var res struct {
		Folders struct {
			Nodes []*Folder `json:"nodes"`
		} `json:"folders"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, xerrors.Errorf("failed to ki.getFolders: %w", err)
	}
	return res.Folders.Nodes, nil
}

func (ki *Kibela) fetchFolders() (map[string]ID, error) {
	ki.foldersOnce.Do(func() {
		if ki.folders != nil {
			return
		}
		folders, err := ki.getFolders()
		if err != nil {
			ki.foldersErr = xerrors.Errorf("failed to ki.setFolders: %w", err)
			return
		}
		folderMap := make(map[string]ID, len(folders))
		for _, fo := range folders {
			folderMap[fo.FullName] = fo.ID
		}
		ki.folders = folderMap
	})
	return ki.folders, ki.foldersErr
}

func (ki *Kibela) fetchFolderID(name string) (ID, error) {
	folders, err := ki.fetchFolders()
	if err != nil {
		return "", xerrors.Errorf("failed to fetchFolderID while setFolderID: %w", err)
	}
	return folders[name], nil
}
