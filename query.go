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

func updateNoteMutation() string {
	return `mutation($id: ID!, $baseNote: NoteInput!, $newNote: NoteInput!) {
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
}

type noteInput struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	GroupIDs  []ID   `json:"groupIds"`
	Folder    string `json:"folderName,omitempty"`
	CoEditing bool   `json:"coediting"`
}
