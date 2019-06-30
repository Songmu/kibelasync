package kibela

import (
	"os"
	"reflect"
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

func newTestMD() *md {
	return &md{
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
}

func TestMD_fullContent(t *testing.T) {
	m := newTestMD()
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

func TestLoadMD(t *testing.T) {
	fpath := "testdata/notes/366.md"
	fi, err := os.Stat(fpath)
	if err != nil {
		t.Fatal(err)
	}
	m, err := loadMD(fpath)
	if err != nil {
		t.Errorf("error should be nil but: %s", err)
	}
	expect := newTestMD()
	expect.filepath = fpath
	expect.UpdatedAt = fi.ModTime()
	if !reflect.DeepEqual(*m, *expect) {
		t.Errorf("got: %+v\nexpect: %+v", *m, *expect)
	}
}
