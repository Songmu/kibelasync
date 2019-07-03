package kibela

import "fmt"

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
	query := fmt.Sprintf("first: %d", num)
	if cursor != "" {
		// cursor is base64 encoded number. ex. "Nw" = 7
		query = fmt.Sprintf(`%s, after: "%s"`, query, cursor)
	}
	return fmt.Sprintf(`{
  notes(%s){
    edges {
      node {
        id
        updatedAt
      }
      cursor
    }
  }
}`, query)
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
