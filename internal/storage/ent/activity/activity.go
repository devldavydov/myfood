// Code generated by ent, DO NOT EDIT.

package activity

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the activity type in the database.
	Label = "activity"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldUserid holds the string denoting the userid field in the database.
	FieldUserid = "userid"
	// FieldTimestamp holds the string denoting the timestamp field in the database.
	FieldTimestamp = "timestamp"
	// FieldActiveCal holds the string denoting the active_cal field in the database.
	FieldActiveCal = "active_cal"
	// Table holds the table name of the activity in the database.
	Table = "activities"
)

// Columns holds all SQL columns for activity fields.
var Columns = []string{
	FieldID,
	FieldUserid,
	FieldTimestamp,
	FieldActiveCal,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the Activity queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByUserid orders the results by the userid field.
func ByUserid(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUserid, opts...).ToFunc()
}

// ByTimestamp orders the results by the timestamp field.
func ByTimestamp(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTimestamp, opts...).ToFunc()
}

// ByActiveCal orders the results by the active_cal field.
func ByActiveCal(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldActiveCal, opts...).ToFunc()
}