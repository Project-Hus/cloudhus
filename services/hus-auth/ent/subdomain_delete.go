// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"hus-auth/ent/predicate"
	"hus-auth/ent/subdomain"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// SubdomainDelete is the builder for deleting a Subdomain entity.
type SubdomainDelete struct {
	config
	hooks    []Hook
	mutation *SubdomainMutation
}

// Where appends a list predicates to the SubdomainDelete builder.
func (sd *SubdomainDelete) Where(ps ...predicate.Subdomain) *SubdomainDelete {
	sd.mutation.Where(ps...)
	return sd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (sd *SubdomainDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, SubdomainMutation](ctx, sd.sqlExec, sd.mutation, sd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (sd *SubdomainDelete) ExecX(ctx context.Context) int {
	n, err := sd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (sd *SubdomainDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(subdomain.Table, sqlgraph.NewFieldSpec(subdomain.FieldID, field.TypeInt))
	if ps := sd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, sd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	sd.mutation.done = true
	return affected, err
}

// SubdomainDeleteOne is the builder for deleting a single Subdomain entity.
type SubdomainDeleteOne struct {
	sd *SubdomainDelete
}

// Where appends a list predicates to the SubdomainDelete builder.
func (sdo *SubdomainDeleteOne) Where(ps ...predicate.Subdomain) *SubdomainDeleteOne {
	sdo.sd.mutation.Where(ps...)
	return sdo
}

// Exec executes the deletion query.
func (sdo *SubdomainDeleteOne) Exec(ctx context.Context) error {
	n, err := sdo.sd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{subdomain.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (sdo *SubdomainDeleteOne) ExecX(ctx context.Context) {
	if err := sdo.Exec(ctx); err != nil {
		panic(err)
	}
}