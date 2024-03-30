// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/devldavydov/myfood/internal/storage/ent/predicate"
	"github.com/devldavydov/myfood/internal/storage/ent/usersettings"
)

// UserSettingsDelete is the builder for deleting a UserSettings entity.
type UserSettingsDelete struct {
	config
	hooks    []Hook
	mutation *UserSettingsMutation
}

// Where appends a list predicates to the UserSettingsDelete builder.
func (usd *UserSettingsDelete) Where(ps ...predicate.UserSettings) *UserSettingsDelete {
	usd.mutation.Where(ps...)
	return usd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (usd *UserSettingsDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, usd.sqlExec, usd.mutation, usd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (usd *UserSettingsDelete) ExecX(ctx context.Context) int {
	n, err := usd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (usd *UserSettingsDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(usersettings.Table, sqlgraph.NewFieldSpec(usersettings.FieldID, field.TypeInt))
	if ps := usd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, usd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	usd.mutation.done = true
	return affected, err
}

// UserSettingsDeleteOne is the builder for deleting a single UserSettings entity.
type UserSettingsDeleteOne struct {
	usd *UserSettingsDelete
}

// Where appends a list predicates to the UserSettingsDelete builder.
func (usdo *UserSettingsDeleteOne) Where(ps ...predicate.UserSettings) *UserSettingsDeleteOne {
	usdo.usd.mutation.Where(ps...)
	return usdo
}

// Exec executes the deletion query.
func (usdo *UserSettingsDeleteOne) Exec(ctx context.Context) error {
	n, err := usdo.usd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{usersettings.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (usdo *UserSettingsDeleteOne) ExecX(ctx context.Context) {
	if err := usdo.Exec(ctx); err != nil {
		panic(err)
	}
}
