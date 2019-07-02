package kibela

import "fmt"

// .data/notes.totalCount
const totalCountQuery = `{
  notes() {
    totalCount
  }
}`

// .data.notes.nodes[]
func listNoteQuery(num int) string {
	return fmt.Sprintf(`{
  notes(first: %d) {
    nodes {
      id
      updatedAt
    }
  }
}`, num)
}

// .data.note
func getNoteQuery(id ID) string {
	return fmt.Sprintf(`{
  note(id: "%s") {
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
    updatedAt
  }
}`, string(id))
}

/*
{
  "data": {
    "notes": {
      "edges": [
        {
          "node": {
            "id": "QmxvZy8zNzA",
            "updatedAt": "2019-06-23T17:39:47.433+09:00"
          },
          "cursor": "NA"
        },
        {
          "node": {
            "id": "QmxvZy8zNzE",
            "updatedAt": "2019-06-23T17:40:09.969+09:00"
          },
          "cursor": "NQ"
        },
        {
          "node": {
            "id": "QmxvZy8zNjg",
            "updatedAt": "2019-06-23T17:39:41.751+09:00"
          },
          "cursor": "Ng"
        }
      ]
    }
  }
}
*/
func listNotePaginateQuery(num int, cursor string) string {
	// cursor is base64 encoded number. ex. "Nw" = 7
	return fmt.Sprintf(`{
  notes(first: %d, after: "%s"){
    edges {
      node {
        id
        updatedAt
      }
      cursor
    }
  }
}`, num, cursor)
}

const totalGroupCountQuery = `{
  groups() {
    totalCount
  }
}`

func listGroupQuery(num int) string {
	return fmt.Sprintf(`{
  groups(first: %d) {
    nodes {
      id
      name
    }
  }
}`, num)
}

const updateNoteMutation = `mutation($id: ID!, $baseNote: NoteInput!, $newNote: NoteInput!) {
  updateNote(input: {
    id: $id,
    baseNote: $baseNote,
    newNote: $newNote,
    draft: false })
  {
    note {
      updatedAt
    }
  }
}`

type noteInput struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	GroupIDs  []ID   `json:"groupIds"`
	Folder    string `json:"folderName,omitempty"`
	CoEditing bool   `json:"coediting"`
}
