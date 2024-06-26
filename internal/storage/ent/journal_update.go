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
	"github.com/devldavydov/myfood/internal/storage/ent/food"
	"github.com/devldavydov/myfood/internal/storage/ent/journal"
	"github.com/devldavydov/myfood/internal/storage/ent/predicate"
)

// JournalUpdate is the builder for updating Journal entities.
type JournalUpdate struct {
	config
	hooks     []Hook
	mutation  *JournalMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the JournalUpdate builder.
func (ju *JournalUpdate) Where(ps ...predicate.Journal) *JournalUpdate {
	ju.mutation.Where(ps...)
	return ju
}

// SetUserid sets the "userid" field.
func (ju *JournalUpdate) SetUserid(i int64) *JournalUpdate {
	ju.mutation.ResetUserid()
	ju.mutation.SetUserid(i)
	return ju
}

// SetNillableUserid sets the "userid" field if the given value is not nil.
func (ju *JournalUpdate) SetNillableUserid(i *int64) *JournalUpdate {
	if i != nil {
		ju.SetUserid(*i)
	}
	return ju
}

// AddUserid adds i to the "userid" field.
func (ju *JournalUpdate) AddUserid(i int64) *JournalUpdate {
	ju.mutation.AddUserid(i)
	return ju
}

// SetTimestamp sets the "timestamp" field.
func (ju *JournalUpdate) SetTimestamp(t time.Time) *JournalUpdate {
	ju.mutation.SetTimestamp(t)
	return ju
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (ju *JournalUpdate) SetNillableTimestamp(t *time.Time) *JournalUpdate {
	if t != nil {
		ju.SetTimestamp(*t)
	}
	return ju
}

// SetMeal sets the "meal" field.
func (ju *JournalUpdate) SetMeal(i int64) *JournalUpdate {
	ju.mutation.ResetMeal()
	ju.mutation.SetMeal(i)
	return ju
}

// SetNillableMeal sets the "meal" field if the given value is not nil.
func (ju *JournalUpdate) SetNillableMeal(i *int64) *JournalUpdate {
	if i != nil {
		ju.SetMeal(*i)
	}
	return ju
}

// AddMeal adds i to the "meal" field.
func (ju *JournalUpdate) AddMeal(i int64) *JournalUpdate {
	ju.mutation.AddMeal(i)
	return ju
}

// SetFoodweight sets the "foodweight" field.
func (ju *JournalUpdate) SetFoodweight(f float64) *JournalUpdate {
	ju.mutation.ResetFoodweight()
	ju.mutation.SetFoodweight(f)
	return ju
}

// SetNillableFoodweight sets the "foodweight" field if the given value is not nil.
func (ju *JournalUpdate) SetNillableFoodweight(f *float64) *JournalUpdate {
	if f != nil {
		ju.SetFoodweight(*f)
	}
	return ju
}

// AddFoodweight adds f to the "foodweight" field.
func (ju *JournalUpdate) AddFoodweight(f float64) *JournalUpdate {
	ju.mutation.AddFoodweight(f)
	return ju
}

// SetFoodID sets the "food" edge to the Food entity by ID.
func (ju *JournalUpdate) SetFoodID(id int) *JournalUpdate {
	ju.mutation.SetFoodID(id)
	return ju
}

// SetFood sets the "food" edge to the Food entity.
func (ju *JournalUpdate) SetFood(f *Food) *JournalUpdate {
	return ju.SetFoodID(f.ID)
}

// Mutation returns the JournalMutation object of the builder.
func (ju *JournalUpdate) Mutation() *JournalMutation {
	return ju.mutation
}

// ClearFood clears the "food" edge to the Food entity.
func (ju *JournalUpdate) ClearFood() *JournalUpdate {
	ju.mutation.ClearFood()
	return ju
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ju *JournalUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, ju.sqlSave, ju.mutation, ju.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ju *JournalUpdate) SaveX(ctx context.Context) int {
	affected, err := ju.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ju *JournalUpdate) Exec(ctx context.Context) error {
	_, err := ju.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ju *JournalUpdate) ExecX(ctx context.Context) {
	if err := ju.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ju *JournalUpdate) check() error {
	if _, ok := ju.mutation.FoodID(); ju.mutation.FoodCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Journal.food"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ju *JournalUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *JournalUpdate {
	ju.modifiers = append(ju.modifiers, modifiers...)
	return ju
}

func (ju *JournalUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := ju.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(journal.Table, journal.Columns, sqlgraph.NewFieldSpec(journal.FieldID, field.TypeInt))
	if ps := ju.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ju.mutation.Userid(); ok {
		_spec.SetField(journal.FieldUserid, field.TypeInt64, value)
	}
	if value, ok := ju.mutation.AddedUserid(); ok {
		_spec.AddField(journal.FieldUserid, field.TypeInt64, value)
	}
	if value, ok := ju.mutation.Timestamp(); ok {
		_spec.SetField(journal.FieldTimestamp, field.TypeTime, value)
	}
	if value, ok := ju.mutation.Meal(); ok {
		_spec.SetField(journal.FieldMeal, field.TypeInt64, value)
	}
	if value, ok := ju.mutation.AddedMeal(); ok {
		_spec.AddField(journal.FieldMeal, field.TypeInt64, value)
	}
	if value, ok := ju.mutation.Foodweight(); ok {
		_spec.SetField(journal.FieldFoodweight, field.TypeFloat64, value)
	}
	if value, ok := ju.mutation.AddedFoodweight(); ok {
		_spec.AddField(journal.FieldFoodweight, field.TypeFloat64, value)
	}
	if ju.mutation.FoodCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   journal.FoodTable,
			Columns: []string{journal.FoodColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(food.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ju.mutation.FoodIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   journal.FoodTable,
			Columns: []string{journal.FoodColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(food.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(ju.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, ju.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{journal.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ju.mutation.done = true
	return n, nil
}

// JournalUpdateOne is the builder for updating a single Journal entity.
type JournalUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *JournalMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUserid sets the "userid" field.
func (juo *JournalUpdateOne) SetUserid(i int64) *JournalUpdateOne {
	juo.mutation.ResetUserid()
	juo.mutation.SetUserid(i)
	return juo
}

// SetNillableUserid sets the "userid" field if the given value is not nil.
func (juo *JournalUpdateOne) SetNillableUserid(i *int64) *JournalUpdateOne {
	if i != nil {
		juo.SetUserid(*i)
	}
	return juo
}

// AddUserid adds i to the "userid" field.
func (juo *JournalUpdateOne) AddUserid(i int64) *JournalUpdateOne {
	juo.mutation.AddUserid(i)
	return juo
}

// SetTimestamp sets the "timestamp" field.
func (juo *JournalUpdateOne) SetTimestamp(t time.Time) *JournalUpdateOne {
	juo.mutation.SetTimestamp(t)
	return juo
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (juo *JournalUpdateOne) SetNillableTimestamp(t *time.Time) *JournalUpdateOne {
	if t != nil {
		juo.SetTimestamp(*t)
	}
	return juo
}

// SetMeal sets the "meal" field.
func (juo *JournalUpdateOne) SetMeal(i int64) *JournalUpdateOne {
	juo.mutation.ResetMeal()
	juo.mutation.SetMeal(i)
	return juo
}

// SetNillableMeal sets the "meal" field if the given value is not nil.
func (juo *JournalUpdateOne) SetNillableMeal(i *int64) *JournalUpdateOne {
	if i != nil {
		juo.SetMeal(*i)
	}
	return juo
}

// AddMeal adds i to the "meal" field.
func (juo *JournalUpdateOne) AddMeal(i int64) *JournalUpdateOne {
	juo.mutation.AddMeal(i)
	return juo
}

// SetFoodweight sets the "foodweight" field.
func (juo *JournalUpdateOne) SetFoodweight(f float64) *JournalUpdateOne {
	juo.mutation.ResetFoodweight()
	juo.mutation.SetFoodweight(f)
	return juo
}

// SetNillableFoodweight sets the "foodweight" field if the given value is not nil.
func (juo *JournalUpdateOne) SetNillableFoodweight(f *float64) *JournalUpdateOne {
	if f != nil {
		juo.SetFoodweight(*f)
	}
	return juo
}

// AddFoodweight adds f to the "foodweight" field.
func (juo *JournalUpdateOne) AddFoodweight(f float64) *JournalUpdateOne {
	juo.mutation.AddFoodweight(f)
	return juo
}

// SetFoodID sets the "food" edge to the Food entity by ID.
func (juo *JournalUpdateOne) SetFoodID(id int) *JournalUpdateOne {
	juo.mutation.SetFoodID(id)
	return juo
}

// SetFood sets the "food" edge to the Food entity.
func (juo *JournalUpdateOne) SetFood(f *Food) *JournalUpdateOne {
	return juo.SetFoodID(f.ID)
}

// Mutation returns the JournalMutation object of the builder.
func (juo *JournalUpdateOne) Mutation() *JournalMutation {
	return juo.mutation
}

// ClearFood clears the "food" edge to the Food entity.
func (juo *JournalUpdateOne) ClearFood() *JournalUpdateOne {
	juo.mutation.ClearFood()
	return juo
}

// Where appends a list predicates to the JournalUpdate builder.
func (juo *JournalUpdateOne) Where(ps ...predicate.Journal) *JournalUpdateOne {
	juo.mutation.Where(ps...)
	return juo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (juo *JournalUpdateOne) Select(field string, fields ...string) *JournalUpdateOne {
	juo.fields = append([]string{field}, fields...)
	return juo
}

// Save executes the query and returns the updated Journal entity.
func (juo *JournalUpdateOne) Save(ctx context.Context) (*Journal, error) {
	return withHooks(ctx, juo.sqlSave, juo.mutation, juo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (juo *JournalUpdateOne) SaveX(ctx context.Context) *Journal {
	node, err := juo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (juo *JournalUpdateOne) Exec(ctx context.Context) error {
	_, err := juo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (juo *JournalUpdateOne) ExecX(ctx context.Context) {
	if err := juo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (juo *JournalUpdateOne) check() error {
	if _, ok := juo.mutation.FoodID(); juo.mutation.FoodCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Journal.food"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (juo *JournalUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *JournalUpdateOne {
	juo.modifiers = append(juo.modifiers, modifiers...)
	return juo
}

func (juo *JournalUpdateOne) sqlSave(ctx context.Context) (_node *Journal, err error) {
	if err := juo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(journal.Table, journal.Columns, sqlgraph.NewFieldSpec(journal.FieldID, field.TypeInt))
	id, ok := juo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Journal.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := juo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, journal.FieldID)
		for _, f := range fields {
			if !journal.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != journal.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := juo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := juo.mutation.Userid(); ok {
		_spec.SetField(journal.FieldUserid, field.TypeInt64, value)
	}
	if value, ok := juo.mutation.AddedUserid(); ok {
		_spec.AddField(journal.FieldUserid, field.TypeInt64, value)
	}
	if value, ok := juo.mutation.Timestamp(); ok {
		_spec.SetField(journal.FieldTimestamp, field.TypeTime, value)
	}
	if value, ok := juo.mutation.Meal(); ok {
		_spec.SetField(journal.FieldMeal, field.TypeInt64, value)
	}
	if value, ok := juo.mutation.AddedMeal(); ok {
		_spec.AddField(journal.FieldMeal, field.TypeInt64, value)
	}
	if value, ok := juo.mutation.Foodweight(); ok {
		_spec.SetField(journal.FieldFoodweight, field.TypeFloat64, value)
	}
	if value, ok := juo.mutation.AddedFoodweight(); ok {
		_spec.AddField(journal.FieldFoodweight, field.TypeFloat64, value)
	}
	if juo.mutation.FoodCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   journal.FoodTable,
			Columns: []string{journal.FoodColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(food.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := juo.mutation.FoodIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   journal.FoodTable,
			Columns: []string{journal.FoodColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(food.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(juo.modifiers...)
	_node = &Journal{config: juo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, juo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{journal.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	juo.mutation.done = true
	return _node, nil
}
