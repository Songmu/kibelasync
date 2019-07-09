package kibela

import (
	"fmt"
	"strings"
)

const totalCountQuery = `{
  notes() {
    totalCount
  }
}`

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

func buildNotesArg(num int, cursor string, hasLimit bool) string {
	var buf = &strings.Builder{}

	fmt.Fprintf(buf, "first: %d", num)
	if cursor != "" {
		// cursor is base64 encoded number. ex. "Nw" = 7
		fmt.Fprintf(buf, ", after: %s", cursor)
	}
	ordering := "PUBLISHED_AT"
	if hasLimit {
		ordering = "CONTENT_UPDATED_AT"
	}
	fmt.Fprintf(buf, ", orderBy: {field: %s, direction: DESC}", ordering)

	// ex. `first: 10, cursor: "Nw", orderBy: {field: PUBLISHED_AT, direction: DESC}`
	return buf.String()
}

func listNoteQuery(num int, hasLimit bool) string {
	return fmt.Sprintf(`{
  notes(%s) {
    nodes {
      id
      updatedAt
    }
  }
}`, buildNotesArg(num, "", hasLimit))
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
func listNotePaginateQuery(num int, cursor string, hasLimit bool) string {
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
}`, buildNotesArg(num, cursor, hasLimit))
}

func listFullNotePaginateQuery(num int, cursor string, hasLimit bool) string {
	return fmt.Sprintf(`{
  notes(%s){
    edges {
      node {
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
        updatedAt
      }
      cursor
    }
  }
}`, buildNotesArg(num, cursor, hasLimit))
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
