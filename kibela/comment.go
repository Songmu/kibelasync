package kibela

import (
	"encoding/json"

	"github.com/Songmu/kibelasync/client"
	"golang.org/x/xerrors"
)

type Comment struct {
	ID      `json:"id"`
	Content string `json:"content"`
	Author  struct {
		Account string `json:"account"`
	}
	PublishedAt Time   `json:"publishedAt"`
	Summary     string `json:"summary"`
}

// GetComment gets kibela comment
func (ki *Kibela) GetComment(num int) (*Comment, error) {
	id := newID(idTypeComment, num)
	data, err := ki.cli.Do(&client.Payload{Query: getCommentQuery(id)})
	if err != nil {
		return nil, xerrors.Errorf("failed to ki.GetComment: %w", err)
	}
	var res struct {
		Comment *Comment `json:"comment"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, xerrors.Errorf("failed to ki.GetComment: %w", err)
	}
	res.Comment.ID = id
	return res.Comment, nil
}
