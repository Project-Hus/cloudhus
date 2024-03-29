// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"hus-auth/ent/connectedsession"
	"hus-auth/ent/hussession"
	"hus-auth/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ConnectedSessionUpdate is the builder for updating ConnectedSession entities.
type ConnectedSessionUpdate struct {
	config
	hooks    []Hook
	mutation *ConnectedSessionMutation
}

// Where appends a list predicates to the ConnectedSessionUpdate builder.
func (csu *ConnectedSessionUpdate) Where(ps ...predicate.ConnectedSession) *ConnectedSessionUpdate {
	csu.mutation.Where(ps...)
	return csu
}

// SetHsid sets the "hsid" field.
func (csu *ConnectedSessionUpdate) SetHsid(u uuid.UUID) *ConnectedSessionUpdate {
	csu.mutation.SetHsid(u)
	return csu
}

// SetService sets the "service" field.
func (csu *ConnectedSessionUpdate) SetService(s string) *ConnectedSessionUpdate {
	csu.mutation.SetService(s)
	return csu
}

// SetCsid sets the "csid" field.
func (csu *ConnectedSessionUpdate) SetCsid(u uuid.UUID) *ConnectedSessionUpdate {
	csu.mutation.SetCsid(u)
	return csu
}

// SetHusSessionID sets the "hus_session" edge to the HusSession entity by ID.
func (csu *ConnectedSessionUpdate) SetHusSessionID(id uuid.UUID) *ConnectedSessionUpdate {
	csu.mutation.SetHusSessionID(id)
	return csu
}

// SetHusSession sets the "hus_session" edge to the HusSession entity.
func (csu *ConnectedSessionUpdate) SetHusSession(h *HusSession) *ConnectedSessionUpdate {
	return csu.SetHusSessionID(h.ID)
}

// Mutation returns the ConnectedSessionMutation object of the builder.
func (csu *ConnectedSessionUpdate) Mutation() *ConnectedSessionMutation {
	return csu.mutation
}

// ClearHusSession clears the "hus_session" edge to the HusSession entity.
func (csu *ConnectedSessionUpdate) ClearHusSession() *ConnectedSessionUpdate {
	csu.mutation.ClearHusSession()
	return csu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (csu *ConnectedSessionUpdate) Save(ctx context.Context) (int, error) {
	return withHooks[int, ConnectedSessionMutation](ctx, csu.sqlSave, csu.mutation, csu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (csu *ConnectedSessionUpdate) SaveX(ctx context.Context) int {
	affected, err := csu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (csu *ConnectedSessionUpdate) Exec(ctx context.Context) error {
	_, err := csu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (csu *ConnectedSessionUpdate) ExecX(ctx context.Context) {
	if err := csu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (csu *ConnectedSessionUpdate) check() error {
	if _, ok := csu.mutation.HusSessionID(); csu.mutation.HusSessionCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "ConnectedSession.hus_session"`)
	}
	return nil
}

func (csu *ConnectedSessionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := csu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(connectedsession.Table, connectedsession.Columns, sqlgraph.NewFieldSpec(connectedsession.FieldID, field.TypeInt))
	if ps := csu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := csu.mutation.Service(); ok {
		_spec.SetField(connectedsession.FieldService, field.TypeString, value)
	}
	if value, ok := csu.mutation.Csid(); ok {
		_spec.SetField(connectedsession.FieldCsid, field.TypeUUID, value)
	}
	if csu.mutation.HusSessionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   connectedsession.HusSessionTable,
			Columns: []string{connectedsession.HusSessionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: hussession.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := csu.mutation.HusSessionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   connectedsession.HusSessionTable,
			Columns: []string{connectedsession.HusSessionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: hussession.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, csu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{connectedsession.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	csu.mutation.done = true
	return n, nil
}

// ConnectedSessionUpdateOne is the builder for updating a single ConnectedSession entity.
type ConnectedSessionUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ConnectedSessionMutation
}

// SetHsid sets the "hsid" field.
func (csuo *ConnectedSessionUpdateOne) SetHsid(u uuid.UUID) *ConnectedSessionUpdateOne {
	csuo.mutation.SetHsid(u)
	return csuo
}

// SetService sets the "service" field.
func (csuo *ConnectedSessionUpdateOne) SetService(s string) *ConnectedSessionUpdateOne {
	csuo.mutation.SetService(s)
	return csuo
}

// SetCsid sets the "csid" field.
func (csuo *ConnectedSessionUpdateOne) SetCsid(u uuid.UUID) *ConnectedSessionUpdateOne {
	csuo.mutation.SetCsid(u)
	return csuo
}

// SetHusSessionID sets the "hus_session" edge to the HusSession entity by ID.
func (csuo *ConnectedSessionUpdateOne) SetHusSessionID(id uuid.UUID) *ConnectedSessionUpdateOne {
	csuo.mutation.SetHusSessionID(id)
	return csuo
}

// SetHusSession sets the "hus_session" edge to the HusSession entity.
func (csuo *ConnectedSessionUpdateOne) SetHusSession(h *HusSession) *ConnectedSessionUpdateOne {
	return csuo.SetHusSessionID(h.ID)
}

// Mutation returns the ConnectedSessionMutation object of the builder.
func (csuo *ConnectedSessionUpdateOne) Mutation() *ConnectedSessionMutation {
	return csuo.mutation
}

// ClearHusSession clears the "hus_session" edge to the HusSession entity.
func (csuo *ConnectedSessionUpdateOne) ClearHusSession() *ConnectedSessionUpdateOne {
	csuo.mutation.ClearHusSession()
	return csuo
}

// Where appends a list predicates to the ConnectedSessionUpdate builder.
func (csuo *ConnectedSessionUpdateOne) Where(ps ...predicate.ConnectedSession) *ConnectedSessionUpdateOne {
	csuo.mutation.Where(ps...)
	return csuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (csuo *ConnectedSessionUpdateOne) Select(field string, fields ...string) *ConnectedSessionUpdateOne {
	csuo.fields = append([]string{field}, fields...)
	return csuo
}

// Save executes the query and returns the updated ConnectedSession entity.
func (csuo *ConnectedSessionUpdateOne) Save(ctx context.Context) (*ConnectedSession, error) {
	return withHooks[*ConnectedSession, ConnectedSessionMutation](ctx, csuo.sqlSave, csuo.mutation, csuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (csuo *ConnectedSessionUpdateOne) SaveX(ctx context.Context) *ConnectedSession {
	node, err := csuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (csuo *ConnectedSessionUpdateOne) Exec(ctx context.Context) error {
	_, err := csuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (csuo *ConnectedSessionUpdateOne) ExecX(ctx context.Context) {
	if err := csuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (csuo *ConnectedSessionUpdateOne) check() error {
	if _, ok := csuo.mutation.HusSessionID(); csuo.mutation.HusSessionCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "ConnectedSession.hus_session"`)
	}
	return nil
}

func (csuo *ConnectedSessionUpdateOne) sqlSave(ctx context.Context) (_node *ConnectedSession, err error) {
	if err := csuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(connectedsession.Table, connectedsession.Columns, sqlgraph.NewFieldSpec(connectedsession.FieldID, field.TypeInt))
	id, ok := csuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "ConnectedSession.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := csuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, connectedsession.FieldID)
		for _, f := range fields {
			if !connectedsession.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != connectedsession.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := csuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := csuo.mutation.Service(); ok {
		_spec.SetField(connectedsession.FieldService, field.TypeString, value)
	}
	if value, ok := csuo.mutation.Csid(); ok {
		_spec.SetField(connectedsession.FieldCsid, field.TypeUUID, value)
	}
	if csuo.mutation.HusSessionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   connectedsession.HusSessionTable,
			Columns: []string{connectedsession.HusSessionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: hussession.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := csuo.mutation.HusSessionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   connectedsession.HusSessionTable,
			Columns: []string{connectedsession.HusSessionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: hussession.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &ConnectedSession{config: csuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, csuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{connectedsession.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	csuo.mutation.done = true
	return _node, nil
}
