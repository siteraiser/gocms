package samplemodels

import "gocms/app"

func UpdatePage(content string) {
	result, err := app.Db.Exec(
		"REPLACE INTO pages (id,link,content) VALUES(?,?,?)",
		1,
		"cms-page",
		content,
	)
	if err != nil {
		panic(err)
	}

	lastid, err := result.LastInsertId()
	if err != nil || lastid < 1 {
		panic(err)
	}
}
