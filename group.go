package kibela

import "fmt"

type Group struct {
	ID   `json:"id"`
	name string `json:"name"`
}

/*
{
  "data": {
    "groups": {
      "totalCount": 2
    }
  }
}
*/
// .data.groups.totalCount
const groupCountQuery = `{
  groups() {
    totalCount
  }
}`

/*
{
  "data": {
    "groups": {
      "nodes": [
        {
          "id": "R3JvdXAvMQ",
          "name": "Home"
        },
        {
          "id": "R3JvdXAvMg",
          "name": "Test"
        }
      ]
    }
  }
}
*/
func groupsQuery() string {
	return fmt.Sprintf(`{
  groups(first:%d) {
    nodes {
      id
      name
    }
  }
  budget {cost}
}`, 100)
}
