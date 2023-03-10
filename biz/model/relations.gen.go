// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameRelation = "relations"

// Relation mapped from table <relations>
type Relation struct {
	RelationID int64 `gorm:"column:relation_id;type:bigint;primaryKey;autoIncrement:true" json:"relation_id"`
	FollowID   int64 `gorm:"column:follow_id;type:bigint;not null" json:"follow_id"`
	FollowerID int64 `gorm:"column:follower_id;type:bigint;not null" json:"follower_id"`
}

// TableName Relation's table name
func (*Relation) TableName() string {
	return TableNameRelation
}
