// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/devldavydov/myfood/internal/storage/ent/usersettings"
)

// UserSettingsCreate is the builder for creating a UserSettings entity.
type UserSettingsCreate struct {
	config
	mutation *UserSettingsMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetUserid sets the "userid" field.
func (usc *UserSettingsCreate) SetUserid(i int64) *UserSettingsCreate {
	usc.mutation.SetUserid(i)
	return usc
}

// SetCalLimit sets the "cal_limit" field.
func (usc *UserSettingsCreate) SetCalLimit(f float64) *UserSettingsCreate {
	usc.mutation.SetCalLimit(f)
	return usc
}

// Mutation returns the UserSettingsMutation object of the builder.
func (usc *UserSettingsCreate) Mutation() *UserSettingsMutation {
	return usc.mutation
}

// Save creates the UserSettings in the database.
func (usc *UserSettingsCreate) Save(ctx context.Context) (*UserSettings, error) {
	return withHooks(ctx, usc.sqlSave, usc.mutation, usc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (usc *UserSettingsCreate) SaveX(ctx context.Context) *UserSettings {
	v, err := usc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (usc *UserSettingsCreate) Exec(ctx context.Context) error {
	_, err := usc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (usc *UserSettingsCreate) ExecX(ctx context.Context) {
	if err := usc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (usc *UserSettingsCreate) check() error {
	if _, ok := usc.mutation.Userid(); !ok {
		return &ValidationError{Name: "userid", err: errors.New(`ent: missing required field "UserSettings.userid"`)}
	}
	if _, ok := usc.mutation.CalLimit(); !ok {
		return &ValidationError{Name: "cal_limit", err: errors.New(`ent: missing required field "UserSettings.cal_limit"`)}
	}
	return nil
}

func (usc *UserSettingsCreate) sqlSave(ctx context.Context) (*UserSettings, error) {
	if err := usc.check(); err != nil {
		return nil, err
	}
	_node, _spec := usc.createSpec()
	if err := sqlgraph.CreateNode(ctx, usc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	usc.mutation.id = &_node.ID
	usc.mutation.done = true
	return _node, nil
}

func (usc *UserSettingsCreate) createSpec() (*UserSettings, *sqlgraph.CreateSpec) {
	var (
		_node = &UserSettings{config: usc.config}
		_spec = sqlgraph.NewCreateSpec(usersettings.Table, sqlgraph.NewFieldSpec(usersettings.FieldID, field.TypeInt))
	)
	_spec.OnConflict = usc.conflict
	if value, ok := usc.mutation.Userid(); ok {
		_spec.SetField(usersettings.FieldUserid, field.TypeInt64, value)
		_node.Userid = value
	}
	if value, ok := usc.mutation.CalLimit(); ok {
		_spec.SetField(usersettings.FieldCalLimit, field.TypeFloat64, value)
		_node.CalLimit = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.UserSettings.Create().
//		SetUserid(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.UserSettingsUpsert) {
//			SetUserid(v+v).
//		}).
//		Exec(ctx)
func (usc *UserSettingsCreate) OnConflict(opts ...sql.ConflictOption) *UserSettingsUpsertOne {
	usc.conflict = opts
	return &UserSettingsUpsertOne{
		create: usc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.UserSettings.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (usc *UserSettingsCreate) OnConflictColumns(columns ...string) *UserSettingsUpsertOne {
	usc.conflict = append(usc.conflict, sql.ConflictColumns(columns...))
	return &UserSettingsUpsertOne{
		create: usc,
	}
}

type (
	// UserSettingsUpsertOne is the builder for "upsert"-ing
	//  one UserSettings node.
	UserSettingsUpsertOne struct {
		create *UserSettingsCreate
	}

	// UserSettingsUpsert is the "OnConflict" setter.
	UserSettingsUpsert struct {
		*sql.UpdateSet
	}
)

// SetUserid sets the "userid" field.
func (u *UserSettingsUpsert) SetUserid(v int64) *UserSettingsUpsert {
	u.Set(usersettings.FieldUserid, v)
	return u
}

// UpdateUserid sets the "userid" field to the value that was provided on create.
func (u *UserSettingsUpsert) UpdateUserid() *UserSettingsUpsert {
	u.SetExcluded(usersettings.FieldUserid)
	return u
}

// AddUserid adds v to the "userid" field.
func (u *UserSettingsUpsert) AddUserid(v int64) *UserSettingsUpsert {
	u.Add(usersettings.FieldUserid, v)
	return u
}

// SetCalLimit sets the "cal_limit" field.
func (u *UserSettingsUpsert) SetCalLimit(v float64) *UserSettingsUpsert {
	u.Set(usersettings.FieldCalLimit, v)
	return u
}

// UpdateCalLimit sets the "cal_limit" field to the value that was provided on create.
func (u *UserSettingsUpsert) UpdateCalLimit() *UserSettingsUpsert {
	u.SetExcluded(usersettings.FieldCalLimit)
	return u
}

// AddCalLimit adds v to the "cal_limit" field.
func (u *UserSettingsUpsert) AddCalLimit(v float64) *UserSettingsUpsert {
	u.Add(usersettings.FieldCalLimit, v)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.UserSettings.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *UserSettingsUpsertOne) UpdateNewValues() *UserSettingsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.UserSettings.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *UserSettingsUpsertOne) Ignore() *UserSettingsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *UserSettingsUpsertOne) DoNothing() *UserSettingsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the UserSettingsCreate.OnConflict
// documentation for more info.
func (u *UserSettingsUpsertOne) Update(set func(*UserSettingsUpsert)) *UserSettingsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&UserSettingsUpsert{UpdateSet: update})
	}))
	return u
}

// SetUserid sets the "userid" field.
func (u *UserSettingsUpsertOne) SetUserid(v int64) *UserSettingsUpsertOne {
	return u.Update(func(s *UserSettingsUpsert) {
		s.SetUserid(v)
	})
}

// AddUserid adds v to the "userid" field.
func (u *UserSettingsUpsertOne) AddUserid(v int64) *UserSettingsUpsertOne {
	return u.Update(func(s *UserSettingsUpsert) {
		s.AddUserid(v)
	})
}

// UpdateUserid sets the "userid" field to the value that was provided on create.
func (u *UserSettingsUpsertOne) UpdateUserid() *UserSettingsUpsertOne {
	return u.Update(func(s *UserSettingsUpsert) {
		s.UpdateUserid()
	})
}

// SetCalLimit sets the "cal_limit" field.
func (u *UserSettingsUpsertOne) SetCalLimit(v float64) *UserSettingsUpsertOne {
	return u.Update(func(s *UserSettingsUpsert) {
		s.SetCalLimit(v)
	})
}

// AddCalLimit adds v to the "cal_limit" field.
func (u *UserSettingsUpsertOne) AddCalLimit(v float64) *UserSettingsUpsertOne {
	return u.Update(func(s *UserSettingsUpsert) {
		s.AddCalLimit(v)
	})
}

// UpdateCalLimit sets the "cal_limit" field to the value that was provided on create.
func (u *UserSettingsUpsertOne) UpdateCalLimit() *UserSettingsUpsertOne {
	return u.Update(func(s *UserSettingsUpsert) {
		s.UpdateCalLimit()
	})
}

// Exec executes the query.
func (u *UserSettingsUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for UserSettingsCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *UserSettingsUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *UserSettingsUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *UserSettingsUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// UserSettingsCreateBulk is the builder for creating many UserSettings entities in bulk.
type UserSettingsCreateBulk struct {
	config
	err      error
	builders []*UserSettingsCreate
	conflict []sql.ConflictOption
}

// Save creates the UserSettings entities in the database.
func (uscb *UserSettingsCreateBulk) Save(ctx context.Context) ([]*UserSettings, error) {
	if uscb.err != nil {
		return nil, uscb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(uscb.builders))
	nodes := make([]*UserSettings, len(uscb.builders))
	mutators := make([]Mutator, len(uscb.builders))
	for i := range uscb.builders {
		func(i int, root context.Context) {
			builder := uscb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*UserSettingsMutation)
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
					_, err = mutators[i+1].Mutate(root, uscb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = uscb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, uscb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, uscb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (uscb *UserSettingsCreateBulk) SaveX(ctx context.Context) []*UserSettings {
	v, err := uscb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (uscb *UserSettingsCreateBulk) Exec(ctx context.Context) error {
	_, err := uscb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uscb *UserSettingsCreateBulk) ExecX(ctx context.Context) {
	if err := uscb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.UserSettings.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.UserSettingsUpsert) {
//			SetUserid(v+v).
//		}).
//		Exec(ctx)
func (uscb *UserSettingsCreateBulk) OnConflict(opts ...sql.ConflictOption) *UserSettingsUpsertBulk {
	uscb.conflict = opts
	return &UserSettingsUpsertBulk{
		create: uscb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.UserSettings.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (uscb *UserSettingsCreateBulk) OnConflictColumns(columns ...string) *UserSettingsUpsertBulk {
	uscb.conflict = append(uscb.conflict, sql.ConflictColumns(columns...))
	return &UserSettingsUpsertBulk{
		create: uscb,
	}
}

// UserSettingsUpsertBulk is the builder for "upsert"-ing
// a bulk of UserSettings nodes.
type UserSettingsUpsertBulk struct {
	create *UserSettingsCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.UserSettings.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *UserSettingsUpsertBulk) UpdateNewValues() *UserSettingsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.UserSettings.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *UserSettingsUpsertBulk) Ignore() *UserSettingsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *UserSettingsUpsertBulk) DoNothing() *UserSettingsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the UserSettingsCreateBulk.OnConflict
// documentation for more info.
func (u *UserSettingsUpsertBulk) Update(set func(*UserSettingsUpsert)) *UserSettingsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&UserSettingsUpsert{UpdateSet: update})
	}))
	return u
}

// SetUserid sets the "userid" field.
func (u *UserSettingsUpsertBulk) SetUserid(v int64) *UserSettingsUpsertBulk {
	return u.Update(func(s *UserSettingsUpsert) {
		s.SetUserid(v)
	})
}

// AddUserid adds v to the "userid" field.
func (u *UserSettingsUpsertBulk) AddUserid(v int64) *UserSettingsUpsertBulk {
	return u.Update(func(s *UserSettingsUpsert) {
		s.AddUserid(v)
	})
}

// UpdateUserid sets the "userid" field to the value that was provided on create.
func (u *UserSettingsUpsertBulk) UpdateUserid() *UserSettingsUpsertBulk {
	return u.Update(func(s *UserSettingsUpsert) {
		s.UpdateUserid()
	})
}

// SetCalLimit sets the "cal_limit" field.
func (u *UserSettingsUpsertBulk) SetCalLimit(v float64) *UserSettingsUpsertBulk {
	return u.Update(func(s *UserSettingsUpsert) {
		s.SetCalLimit(v)
	})
}

// AddCalLimit adds v to the "cal_limit" field.
func (u *UserSettingsUpsertBulk) AddCalLimit(v float64) *UserSettingsUpsertBulk {
	return u.Update(func(s *UserSettingsUpsert) {
		s.AddCalLimit(v)
	})
}

// UpdateCalLimit sets the "cal_limit" field to the value that was provided on create.
func (u *UserSettingsUpsertBulk) UpdateCalLimit() *UserSettingsUpsertBulk {
	return u.Update(func(s *UserSettingsUpsert) {
		s.UpdateCalLimit()
	})
}

// Exec executes the query.
func (u *UserSettingsUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the UserSettingsCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for UserSettingsCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *UserSettingsUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
