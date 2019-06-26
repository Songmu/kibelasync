package kibela

import "fmt"

/*
   {
     "id": "QmxvZy8zNjY",
     "title": "APIテストpublic",
     "content": "コンテント!\nコンテント",
     "coediting": true,
     "folderName": "testtop/testsub1",
     "groups": [
       {
         "name": "Home",
         "id": "R3JvdXAvMQ"
       }
     ],
     "author": {
       "account": "Songmu"
     },
     "createdAt": "2019-06-23T16:54:09.447+09:00",
     "publishedAt": "2019-06-23T16:54:09.444+09:00",
     "contentUpdatedAt": "2019-06-23T16:54:09.445+09:00",
     "updatedAt": "2019-06-23T17:22:38.496+09:00"
   },
*/
type note struct {
	ID        `json:"id"`
	Title     string `json:"title"`
	CoEditing bool   `json:"coediting"`
	Folder    string `json:"folderName"`
	Groups    struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	Author struct {
		Account string `json:"account"`
	}
	CreatedAt        Time `json:"createdAt"`
	PublishedAt      Time `json:"publishedAt"`
	UpdatedAt        Time `json:"updatedAt"`
	ContentUpdatedAt Time `json:"contentUpdatedAt"`
}

// .data/notes.totalCount
const totalCountQuery = `{
  notes() {
    totalCount
  }
}`

// .data.notes.nodes[]
func listNoteQuery() string {
	return fmt.Sprintf(`{
  notes(first: %d) {
    nodes {
      id
      updatedAt
    }
  }
}`, 100)
}

// .data.note
func getNoteQuery(id string) string {
	return fmt.Sprintf(`{
  note(id: "%s") {
    id
    title
    content
    coediting
    folderName
    groups {
      name
      id
    }
    author {
      account
    }
    createdAt
    publishedAt
    contentUpdatedAt
    updatedAt
  }
}`, id)
}
