// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/devldavydov/myfood/internal/storage/ent/predicate"
	"github.com/devldavydov/myfood/internal/storage/ent/weight"
)

// WeightDelete is the builder for deleting a Weight entity.
type WeightDelete struct {
	config
	hooks    []Hook
	mutation *WeightMutation
}

// Where appends a list predicates to the WeightDelete builder.
func (wd *WeightDelete) Where(ps ...predicate.Weight) *WeightDelete {
	wd.mutation.Where(ps...)
	return wd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (wd *WeightDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, wd.sqlExec, wd.mutation, wd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (wd *WeightDelete) ExecX(ctx context.Context) int {
	n, err := wd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (wd *WeightDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(weight.Table, sqlgraph.NewFieldSpec(weight.FieldID, field.TypeInt))
	if ps := wd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, wd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	wd.mutation.done = true
	return affected, err
}

// WeightDeleteOne is the builder for deleting a single Weight entity.
type WeightDeleteOne struct {
	wd *WeightDelete
}

// Where appends a list predicates to the WeightDelete builder.
func (wdo *WeightDeleteOne) Where(ps ...predicate.Weight) *WeightDeleteOne {
	wdo.wd.mutation.Where(ps...)
	return wdo
}

// Exec executes the deletion query.
func (wdo *WeightDeleteOne) Exec(ctx context.Context) error {
	n, err := wdo.wd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{weight.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (wdo *WeightDeleteOne) ExecX(ctx context.Context) {
	if err := wdo.Exec(ctx); err != nil {
		panic(err)
	}
}