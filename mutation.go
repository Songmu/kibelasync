package kibelasync

type noteInput struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	GroupIDs  []string `json:"groupIds"`
	Folder    string   `json:"folderName,omitempty"`
	CoEditing bool     `json:"coediting"`
}

const updateNoteMutation = `mutation($id: ID!, $baseNote: NoteInput!, $newNote: NoteInput!) {
  updateNote(input: {
    id: $id,
    baseNote: $baseNote,
    newNote: $newNote,
    draft: false })
  {
    note {
      author {
        account
      }
      updatedAt
    }
  }
}`

const createNoteMutation = `mutation ($input: CreateNoteInput!) {
  createNote(input: $input) {
    note {
      id
      updatedAt
      groups {
        name
      }
      author {
        account
      }
    }
  }
}`
