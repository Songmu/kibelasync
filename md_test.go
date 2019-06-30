package kibela

import (
	"testing"
)

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
     "createdAt": "2019-06-23T16:54:09.447Z",
     "publishedAt": "2019-06-23T16:54:09.444Z",
     "contentUpdatedAt": "2019-06-23T16:54:09.445Z",
     "updatedAt": "2019-06-23T17:22:38.496Z"
   },
*/

func TestMD_fullContent(t *testing.T) {
	m := &md{
		ID: ID("QmxvZy8zNjY"),
		FrontMatter: &meta{
			Title:     "たいとる！",
			CoEditing: true,
			Folder:    "hoge/fuga",
			Groups:    []string{"Public", "Hobby"},
			Author:    "Songmu",
		},
		Content: "Hello World!\nこんにちは!\n",
	}

	out := m.fullContent()
	expect := `---
author: Songmu
coediting: true
folder: hoge/fuga
groups:
- Public
- Hobby
title: たいとる！
---

Hello World!
こんにちは!
`
	if out != expect {
		t.Errorf("m.fullContent() = got:\n%s\nexpect:\n%s\n", out, expect)
	}
}
