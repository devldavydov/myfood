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
)

// JournalCreate is the builder for creating a Journal entity.
type JournalCreate struct {
	config
	mutation *JournalMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetUserid sets the "userid" field.
func (jc *JournalCreate) SetUserid(i int64) *JournalCreate {
	jc.mutation.SetUserid(i)
	return jc
}

// SetTimestamp sets the "timestamp" field.
func (jc *JournalCreate) SetTimestamp(t time.Time) *JournalCreate {
	jc.mutation.SetTimestamp(t)
	return jc
}

// SetMeal sets the "meal" field.
func (jc *JournalCreate) SetMeal(i int64) *JournalCreate {
	jc.mutation.SetMeal(i)
	return jc
}

// SetFoodweight sets the "foodweight" field.
func (jc *JournalCreate) SetFoodweight(f float64) *JournalCreate {
	jc.mutation.SetFoodweight(f)
	return jc
}

// SetFoodID sets the "food" edge to the Food entity by ID.
func (jc *JournalCreate) SetFoodID(id int) *JournalCreate {
	jc.mutation.SetFoodID(id)
	return jc
}

// SetFood sets the "food" edge to the Food entity.
func (jc *JournalCreate) SetFood(f *Food) *JournalCreate {
	return jc.SetFoodID(f.ID)
}

// Mutation returns the JournalMutation object of the builder.
func (jc *JournalCreate) Mutation() *JournalMutation {
	return jc.mutation
}

// Save creates the Journal in the database.
func (jc *JournalCreate) Save(ctx context.Context) (*Journal, error) {
	return withHooks(ctx, jc.sqlSave, jc.mutation, jc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (jc *JournalCreate) SaveX(ctx context.Context) *Journal {
	v, err := jc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (jc *JournalCreate) Exec(ctx context.Context) error {
	_, err := jc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (jc *JournalCreate) ExecX(ctx context.Context) {
	if err := jc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (jc *JournalCreate) check() error {
	if _, ok := jc.mutation.Userid(); !ok {
		return &ValidationError{Name: "userid", err: errors.New(`ent: missing required field "Journal.userid"`)}
	}
	if _, ok := jc.mutation.Timestamp(); !ok {
		return &ValidationError{Name: "timestamp", err: errors.New(`ent: missing required field "Journal.timestamp"`)}
	}
	if _, ok := jc.mutation.Meal(); !ok {
		return &ValidationError{Name: "meal", err: errors.New(`ent: missing required field "Journal.meal"`)}
	}
	if _, ok := jc.mutation.Foodweight(); !ok {
		return &ValidationError{Name: "foodweight", err: errors.New(`ent: missing required field "Journal.foodweight"`)}
	}
	if _, ok := jc.mutation.FoodID(); !ok {
		return &ValidationError{Name: "food", err: errors.New(`ent: missing required edge "Journal.food"`)}
	}
	return nil
}

func (jc *JournalCreate) sqlSave(ctx context.Context) (*Journal, error) {
	if err := jc.check(); err != nil {
		return nil, err
	}
	_node, _spec := jc.createSpec()
	if err := sqlgraph.CreateNode(ctx, jc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	jc.mutation.id = &_node.ID
	jc.mutation.done = true
	return _node, nil
}

func (jc *JournalCreate) createSpec() (*Journal, *sqlgraph.CreateSpec) {
	var (
		_node = &Journal{config: jc.config}
		_spec = sqlgraph.NewCreateSpec(journal.Table, sqlgraph.NewFieldSpec(journal.FieldID, field.TypeInt))
	)
	_spec.OnConflict = jc.conflict
	if value, ok := jc.mutation.Userid(); ok {
		_spec.SetField(journal.FieldUserid, field.TypeInt64, value)
		_node.Userid = value
	}
	if value, ok := jc.mutation.Timestamp(); ok {
		_spec.SetField(journal.FieldTimestamp, field.TypeTime, value)
		_node.Timestamp = value
	}
	if value, ok := jc.mutation.Meal(); ok {
		_spec.SetField(journal.FieldMeal, field.TypeInt64, value)
		_node.Meal = value
	}
	if value, ok := jc.mutation.Foodweight(); ok {
		_spec.SetField(journal.FieldFoodweight, field.TypeFloat64, value)
		_node.Foodweight = value
	}
	if nodes := jc.mutation.FoodIDs(); len(nodes) > 0 {
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
		_node.food_journals = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Journal.Create().
//		SetUserid(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.JournalUpsert) {
//			SetUserid(v+v).
//		}).
//		Exec(ctx)
func (jc *JournalCreate) OnConflict(opts ...sql.ConflictOption) *JournalUpsertOne {
	jc.conflict = opts
	return &JournalUpsertOne{
		create: jc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Journal.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (jc *JournalCreate) OnConflictColumns(columns ...string) *JournalUpsertOne {
	jc.conflict = append(jc.conflict, sql.ConflictColumns(columns...))
	return &JournalUpsertOne{
		create: jc,
	}
}

type (
	// JournalUpsertOne is the builder for "upsert"-ing
	//  one Journal node.
	JournalUpsertOne struct {
		create *JournalCreate
	}

	// JournalUpsert is the "OnConflict" setter.
	JournalUpsert struct {
		*sql.UpdateSet
	}
)

// SetUserid sets the "userid" field.
func (u *JournalUpsert) SetUserid(v int64) *JournalUpsert {
	u.Set(journal.FieldUserid, v)
	return u
}

// UpdateUserid sets the "userid" field to the value that was provided on create.
func (u *JournalUpsert) UpdateUserid() *JournalUpsert {
	u.SetExcluded(journal.FieldUserid)
	return u
}

// AddUserid adds v to the "userid" field.
func (u *JournalUpsert) AddUserid(v int64) *JournalUpsert {
	u.Add(journal.FieldUserid, v)
	return u
}

// SetTimestamp sets the "timestamp" field.
func (u *JournalUpsert) SetTimestamp(v time.Time) *JournalUpsert {
	u.Set(journal.FieldTimestamp, v)
	return u
}

// UpdateTimestamp sets the "timestamp" field to the value that was provided on create.
func (u *JournalUpsert) UpdateTimestamp() *JournalUpsert {
	u.SetExcluded(journal.FieldTimestamp)
	return u
}

// SetMeal sets the "meal" field.
func (u *JournalUpsert) SetMeal(v int64) *JournalUpsert {
	u.Set(journal.FieldMeal, v)
	return u
}

// UpdateMeal sets the "meal" field to the value that was provided on create.
func (u *JournalUpsert) UpdateMeal() *JournalUpsert {
	u.SetExcluded(journal.FieldMeal)
	return u
}

// AddMeal adds v to the "meal" field.
func (u *JournalUpsert) AddMeal(v int64) *JournalUpsert {
	u.Add(journal.FieldMeal, v)
	return u
}

// SetFoodweight sets the "foodweight" field.
func (u *JournalUpsert) SetFoodweight(v float64) *JournalUpsert {
	u.Set(journal.FieldFoodweight, v)
	return u
}

// UpdateFoodweight sets the "foodweight" field to the value that was provided on create.
func (u *JournalUpsert) UpdateFoodweight() *JournalUpsert {
	u.SetExcluded(journal.FieldFoodweight)
	return u
}

// AddFoodweight adds v to the "foodweight" field.
func (u *JournalUpsert) AddFoodweight(v float64) *JournalUpsert {
	u.Add(journal.FieldFoodweight, v)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Journal.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *JournalUpsertOne) UpdateNewValues() *JournalUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Journal.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *JournalUpsertOne) Ignore() *JournalUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *JournalUpsertOne) DoNothing() *JournalUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the JournalCreate.OnConflict
// documentation for more info.
func (u *JournalUpsertOne) Update(set func(*JournalUpsert)) *JournalUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&JournalUpsert{UpdateSet: update})
	}))
	return u
}

// SetUserid sets the "userid" field.
func (u *JournalUpsertOne) SetUserid(v int64) *JournalUpsertOne {
	return u.Update(func(s *JournalUpsert) {
		s.SetUserid(v)
	})
}

// AddUserid adds v to the "userid" field.
func (u *JournalUpsertOne) AddUserid(v int64) *JournalUpsertOne {
	return u.Update(func(s *JournalUpsert) {
		s.AddUserid(v)
	})
}

// UpdateUserid sets the "userid" field to the value that was provided on create.
func (u *JournalUpsertOne) UpdateUserid() *JournalUpsertOne {
	return u.Update(func(s *JournalUpsert) {
		s.UpdateUserid()
	})
}

// SetTimestamp sets the "timestamp" field.
func (u *JournalUpsertOne) SetTimestamp(v time.Time) *JournalUpsertOne {
	return u.Update(func(s *JournalUpsert) {
		s.SetTimestamp(v)
	})
}

// UpdateTimestamp sets the "timestamp" field to the value that was provided on create.
func (u *JournalUpsertOne) UpdateTimestamp() *JournalUpsertOne {
	return u.Update(func(s *JournalUpsert) {
		s.UpdateTimestamp()
	})
}

// SetMeal sets the "meal" field.
func (u *JournalUpsertOne) SetMeal(v int64) *JournalUpsertOne {
	return u.Update(func(s *JournalUpsert) {
		s.SetMeal(v)
	})
}

// AddMeal adds v to the "meal" field.
func (u *JournalUpsertOne) AddMeal(v int64) *JournalUpsertOne {
	return u.Update(func(s *JournalUpsert) {
		s.AddMeal(v)
	})
}

// UpdateMeal sets the "meal" field to the value that was provided on create.
func (u *JournalUpsertOne) UpdateMeal() *JournalUpsertOne {
	return u.Update(func(s *JournalUpsert) {
		s.UpdateMeal()
	})
}

// SetFoodweight sets the "foodweight" field.
func (u *JournalUpsertOne) SetFoodweight(v float64) *JournalUpsertOne {
	return u.Update(func(s *JournalUpsert) {
		s.SetFoodweight(v)
	})
}

// AddFoodweight adds v to the "foodweight" field.
func (u *JournalUpsertOne) AddFoodweight(v float64) *JournalUpsertOne {
	return u.Update(func(s *JournalUpsert) {
		s.AddFoodweight(v)
	})
}

// UpdateFoodweight sets the "foodweight" field to the value that was provided on create.
func (u *JournalUpsertOne) UpdateFoodweight() *JournalUpsertOne {
	return u.Update(func(s *JournalUpsert) {
		s.UpdateFoodweight()
	})
}

// Exec executes the query.
func (u *JournalUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for JournalCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *JournalUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *JournalUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *JournalUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// JournalCreateBulk is the builder for creating many Journal entities in bulk.
type JournalCreateBulk struct {
	config
	err      error
	builders []*JournalCreate
	conflict []sql.ConflictOption
}

// Save creates the Journal entities in the database.
func (jcb *JournalCreateBulk) Save(ctx context.Context) ([]*Journal, error) {
	if jcb.err != nil {
		return nil, jcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(jcb.builders))
	nodes := make([]*Journal, len(jcb.builders))
	mutators := make([]Mutator, len(jcb.builders))
	for i := range jcb.builders {
		func(i int, root context.Context) {
			builder := jcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*JournalMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, jcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = jcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, jcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, jcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (jcb *JournalCreateBulk) SaveX(ctx context.Context) []*Journal {
	v, err := jcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (jcb *JournalCreateBulk) Exec(ctx context.Context) error {
	_, err := jcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (jcb *JournalCreateBulk) ExecX(ctx context.Context) {
	if err := jcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Journal.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.JournalUpsert) {
//			SetUserid(v+v).
//		}).
//		Exec(ctx)
func (jcb *JournalCreateBulk) OnConflict(opts ...sql.ConflictOption) *JournalUpsertBulk {
	jcb.conflict = opts
	return &JournalUpsertBulk{
		create: jcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Journal.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (jcb *JournalCreateBulk) OnConflictColumns(columns ...string) *JournalUpsertBulk {
	jcb.conflict = append(jcb.conflict, sql.ConflictColumns(columns...))
	return &JournalUpsertBulk{
		create: jcb,
	}
}

// JournalUpsertBulk is the builder for "upsert"-ing
// a bulk of Journal nodes.
type JournalUpsertBulk struct {
	create *JournalCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Journal.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *JournalUpsertBulk) UpdateNewValues() *JournalUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Journal.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *JournalUpsertBulk) Ignore() *JournalUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *JournalUpsertBulk) DoNothing() *JournalUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the JournalCreateBulk.OnConflict
// documentation for more info.
func (u *JournalUpsertBulk) Update(set func(*JournalUpsert)) *JournalUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&JournalUpsert{UpdateSet: update})
	}))
	return u
}

// SetUserid sets the "userid" field.
func (u *JournalUpsertBulk) SetUserid(v int64) *JournalUpsertBulk {
	return u.Update(func(s *JournalUpsert) {
		s.SetUserid(v)
	})
}

// AddUserid adds v to the "userid" field.
func (u *JournalUpsertBulk) AddUserid(v int64) *JournalUpsertBulk {
	return u.Update(func(s *JournalUpsert) {
		s.AddUserid(v)
	})
}

// UpdateUserid sets the "userid" field to the value that was provided on create.
func (u *JournalUpsertBulk) UpdateUserid() *JournalUpsertBulk {
	return u.Update(func(s *JournalUpsert) {
		s.UpdateUserid()
	})
}

// SetTimestamp sets the "timestamp" field.
func (u *JournalUpsertBulk) SetTimestamp(v time.Time) *JournalUpsertBulk {
	return u.Update(func(s *JournalUpsert) {
		s.SetTimestamp(v)
	})
}

// UpdateTimestamp sets the "timestamp" field to the value that was provided on create.
func (u *JournalUpsertBulk) UpdateTimestamp() *JournalUpsertBulk {
	return u.Update(func(s *JournalUpsert) {
		s.UpdateTimestamp()
	})
}

// SetMeal sets the "meal" field.
func (u *JournalUpsertBulk) SetMeal(v int64) *JournalUpsertBulk {
	return u.Update(func(s *JournalUpsert) {
		s.SetMeal(v)
	})
}

// AddMeal adds v to the "meal" field.
func (u *JournalUpsertBulk) AddMeal(v int64) *JournalUpsertBulk {
	return u.Update(func(s *JournalUpsert) {
		s.AddMeal(v)
	})
}

// UpdateMeal sets the "meal" field to the value that was provided on create.
func (u *JournalUpsertBulk) UpdateMeal() *JournalUpsertBulk {
	return u.Update(func(s *JournalUpsert) {
		s.UpdateMeal()
	})
}

// SetFoodweight sets the "foodweight" field.
func (u *JournalUpsertBulk) SetFoodweight(v float64) *JournalUpsertBulk {
	return u.Update(func(s *JournalUpsert) {
		s.SetFoodweight(v)
	})
}

// AddFoodweight adds v to the "foodweight" field.
func (u *JournalUpsertBulk) AddFoodweight(v float64) *JournalUpsertBulk {
	return u.Update(func(s *JournalUpsert) {
		s.AddFoodweight(v)
	})
}

// UpdateFoodweight sets the "foodweight" field to the value that was provided on create.
func (u *JournalUpsertBulk) UpdateFoodweight() *JournalUpsertBulk {
	return u.Update(func(s *JournalUpsert) {
		s.UpdateFoodweight()
	})
}

// Exec executes the query.
func (u *JournalUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the JournalCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for JournalCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *JournalUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}