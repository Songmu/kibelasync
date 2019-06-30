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
