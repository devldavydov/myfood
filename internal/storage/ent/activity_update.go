// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/devldavydov/myfood/internal/storage/ent/activity"
	"github.com/devldavydov/myfood/internal/storage/ent/predicate"
)

// ActivityUpdate is the builder for updating Activity entities.
type ActivityUpdate struct {
	config
	hooks     []Hook
	mutation  *ActivityMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ActivityUpdate builder.
func (au *ActivityUpdate) Where(ps ...predicate.Activity) *ActivityUpdate {
	au.mutation.Where(ps...)
	return au
}

// SetUserid sets the "userid" field.
func (au *ActivityUpdate) SetUserid(i int64) *ActivityUpdate {
	au.mutation.ResetUserid()
	au.mutation.SetUserid(i)
	return au
}

// SetNillableUserid sets the "userid" field if the given value is not nil.
func (au *ActivityUpdate) SetNillableUserid(i *int64) *ActivityUpdate {
	if i != nil {
		au.SetUserid(*i)
	}
	return au
}

// AddUserid adds i to the "userid" field.
func (au *ActivityUpdate) AddUserid(i int64) *ActivityUpdate {
	au.mutation.AddUserid(i)
	return au
}

// SetTimestamp sets the "timestamp" field.
func (au *ActivityUpdate) SetTimestamp(t time.Time) *ActivityUpdate {
	au.mutation.SetTimestamp(t)
	return au
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (au *ActivityUpdate) SetNillableTimestamp(t *time.Time) *ActivityUpdate {
	if t != nil {
		au.SetTimestamp(*t)
	}
	return au
}

// SetActiveCal sets the "active_cal" field.
func (au *ActivityUpdate) SetActiveCal(f float64) *ActivityUpdate {
	au.mutation.ResetActiveCal()
	au.mutation.SetActiveCal(f)
	return au
}

// SetNillableActiveCal sets the "active_cal" field if the given value is not nil.
func (au *ActivityUpdate) SetNillableActiveCal(f *float64) *ActivityUpdate {
	if f != nil {
		au.SetActiveCal(*f)
	}
	return au
}

// AddActiveCal adds f to the "active_cal" field.
func (au *ActivityUpdate) AddActiveCal(f float64) *ActivityUpdate {
	au.mutation.AddActiveCal(f)
	return au
}

// Mutation returns the ActivityMutation object of the builder.
func (au *ActivityUpdate) Mutation() *ActivityMutation {
	return au.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (au *ActivityUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, au.sqlSave, au.mutation, au.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (au *ActivityUpdate) SaveX(ctx context.Context) int {
	affected, err := au.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (au *ActivityUpdate) Exec(ctx context.Context) error {
	_, err := au.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (au *ActivityUpdate) ExecX(ctx context.Context) {
	if err := au.Exec(ctx); err != nil {
		panic(err)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (au *ActivityUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ActivityUpdate {
	au.modifiers = append(au.modifiers, modifiers...)
	return au
}

func (au *ActivityUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(activity.Table, activity.Columns, sqlgraph.NewFieldSpec(activity.FieldID, field.TypeInt))
	if ps := au.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := au.mutation.Userid(); ok {
		_spec.SetField(activity.FieldUserid, field.TypeInt64, value)
	}
	if value, ok := au.mutation.AddedUserid(); ok {
		_spec.AddField(activity.FieldUserid, field.TypeInt64, value)
	}
	if value, ok := au.mutation.Timestamp(); ok {
		_spec.SetField(activity.FieldTimestamp, field.TypeTime, value)
	}
	if value, ok := au.mutation.ActiveCal(); ok {
		_spec.SetField(activity.FieldActiveCal, field.TypeFloat64, value)
	}
	if value, ok := au.mutation.AddedActiveCal(); ok {
		_spec.AddField(activity.FieldActiveCal, field.TypeFloat64, value)
	}
	_spec.AddModifiers(au.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, au.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{activity.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	au.mutation.done = true
	return n, nil
}

// ActivityUpdateOne is the builder for updating a single Activity entity.
type ActivityUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ActivityMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUserid sets the "userid" field.
func (auo *ActivityUpdateOne) SetUserid(i int64) *ActivityUpdateOne {
	auo.mutation.ResetUserid()
	auo.mutation.SetUserid(i)
	return auo
}

// SetNillableUserid sets the "userid" field if the given value is not nil.
func (auo *ActivityUpdateOne) SetNillableUserid(i *int64) *ActivityUpdateOne {
	if i != nil {
		auo.SetUserid(*i)
	}
	return auo
}

// AddUserid adds i to the "userid" field.
func (auo *ActivityUpdateOne) AddUserid(i int64) *ActivityUpdateOne {
	auo.mutation.AddUserid(i)
	return auo
}

// SetTimestamp sets the "timestamp" field.
func (auo *ActivityUpdateOne) SetTimestamp(t time.Time) *ActivityUpdateOne {
	auo.mutation.SetTimestamp(t)
	return auo
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (auo *ActivityUpdateOne) SetNillableTimestamp(t *time.Time) *ActivityUpdateOne {
	if t != nil {
		auo.SetTimestamp(*t)
	}
	return auo
}

// SetActiveCal sets the "active_cal" field.
func (auo *ActivityUpdateOne) SetActiveCal(f float64) *ActivityUpdateOne {
	auo.mutation.ResetActiveCal()
	auo.mutation.SetActiveCal(f)
	return auo
}

// SetNillableActiveCal sets the "active_cal" field if the given value is not nil.
func (auo *ActivityUpdateOne) SetNillableActiveCal(f *float64) *ActivityUpdateOne {
	if f != nil {
		auo.SetActiveCal(*f)
	}
	return auo
}

// AddActiveCal adds f to the "active_cal" field.
func (auo *ActivityUpdateOne) AddActiveCal(f float64) *ActivityUpdateOne {
	auo.mutation.AddActiveCal(f)
	return auo
}

// Mutation returns the ActivityMutation object of the builder.
func (auo *ActivityUpdateOne) Mutation() *ActivityMutation {
	return auo.mutation
}

// Where appends a list predicates to the ActivityUpdate builder.
func (auo *ActivityUpdateOne) Where(ps ...predicate.Activity) *ActivityUpdateOne {
	auo.mutation.Where(ps...)
	return auo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (auo *ActivityUpdateOne) Select(field string, fields ...string) *ActivityUpdateOne {
	auo.fields = append([]string{field}, fields...)
	return auo
}

// Save executes the query and returns the updated Activity entity.
func (auo *ActivityUpdateOne) Save(ctx context.Context) (*Activity, error) {
	return withHooks(ctx, auo.sqlSave, auo.mutation, auo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (auo *ActivityUpdateOne) SaveX(ctx context.Context) *Activity {
	node, err := auo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (auo *ActivityUpdateOne) Exec(ctx context.Context) error {
	_, err := auo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (auo *ActivityUpdateOne) ExecX(ctx context.Context) {
	if err := auo.Exec(ctx); err != nil {
		panic(err)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (auo *ActivityUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ActivityUpdateOne {
	auo.modifiers = append(auo.modifiers, modifiers...)
	return auo
}

func (auo *ActivityUpdateOne) sqlSave(ctx context.Context) (_node *Activity, err error) {
	_spec := sqlgraph.NewUpdateSpec(activity.Table, activity.Columns, sqlgraph.NewFieldSpec(activity.FieldID, field.TypeInt))
	id, ok := auo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Activity.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := auo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, activity.FieldID)
		for _, f := range fields {
			if !activity.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != activity.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := auo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := auo.mutation.Userid(); ok {
		_spec.SetField(activity.FieldUserid, field.TypeInt64, value)
	}
	if value, ok := auo.mutation.AddedUserid(); ok {
		_spec.AddField(activity.FieldUserid, field.TypeInt64, value)
	}
	if value, ok := auo.mutation.Timestamp(); ok {
		_spec.SetField(activity.FieldTimestamp, field.TypeTime, value)
	}
	if value, ok := auo.mutation.ActiveCal(); ok {
		_spec.SetField(activity.FieldActiveCal, field.TypeFloat64, value)
	}
	if value, ok := auo.mutation.AddedActiveCal(); ok {
		_spec.AddField(activity.FieldActiveCal, field.TypeFloat64, value)
	}
	_spec.AddModifiers(auo.modifiers...)
	_node = &Activity{config: auo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, auo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{activity.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	auo.mutation.done = true
	return _node, nil
}
