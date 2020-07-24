package kibela

import (
	"context"
	"encoding/json"

	"github.com/Songmu/kibelasync/client"
	"golang.org/x/xerrors"
)

// Comment represents comment of Kibela
type Comment struct {
	ID          `json:"id"`
	Content     string `json:"content"`
	Author      User   `json:"author"`
	PublishedAt Time   `json:"publishedAt"`
	Summary     string `json:"summary"`
}

// GetComment gets kibela comment
func (ki *Kibela) GetComment(ctx context.Context, num int) (*Comment, error) {
	id := newID(idTypeComment, num)
	data, err := ki.cli.Do(ctx, &client.Payload{Query: getCommentQuery(id)})
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
