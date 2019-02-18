package module

type CaseFlow struct {
	Id      int64
	CaseKid int64 `xorm:"bigint(20) notnull 'case_kid'"`
}




