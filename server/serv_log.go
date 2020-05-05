package main

import (
	DBCenter "../DataBaseCenter"
	DBModel "../DataBaseCenter/DataBaseModel"
)

func record_visit_log(visitId string, categoryId string) {
	vi_log := DBModel.Visit_log{VisitId: visitId, CategoryId: categoryId}
	DBCenter.InsertTable(vi_log)
}
