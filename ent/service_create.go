// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"hus-auth/ent/service"
	"hus-auth/ent/subdomain"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// ServiceCreate is the builder for creating a Service entity.
type ServiceCreate struct {
	config
	mutation *ServiceMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (sc *ServiceCreate) SetName(s string) *ServiceCreate {
	sc.mutation.SetName(s)
	return sc
}

// SetProtocol sets the "protocol" field.
func (sc *ServiceCreate) SetProtocol(s string) *ServiceCreate {
	sc.mutation.SetProtocol(s)
	return sc
}

// SetDomain sets the "domain" field.
func (sc *ServiceCreate) SetDomain(s string) *ServiceCreate {
	sc.mutation.SetDomain(s)
	return sc
}

// SetCreatedAt sets the "created_at" field.
func (sc *ServiceCreate) SetCreatedAt(t time.Time) *ServiceCreate {
	sc.mutation.SetCreatedAt(t)
	return sc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (sc *ServiceCreate) SetNillableCreatedAt(t *time.Time) *ServiceCreate {
	if t != nil {
		sc.SetCreatedAt(*t)
	}
	return sc
}

// SetUpdatedAt sets the "updated_at" field.
func (sc *ServiceCreate) SetUpdatedAt(t time.Time) *ServiceCreate {
	sc.mutation.SetUpdatedAt(t)
	return sc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (sc *ServiceCreate) SetNillableUpdatedAt(t *time.Time) *ServiceCreate {
	if t != nil {
		sc.SetUpdatedAt(*t)
	}
	return sc
}

// AddSubdomainIDs adds the "subdomains" edge to the Subdomain entity by IDs.
func (sc *ServiceCreate) AddSubdomainIDs(ids ...int) *ServiceCreate {
	sc.mutation.AddSubdomainIDs(ids...)
	return sc
}

// AddSubdomains adds the "subdomains" edges to the Subdomain entity.
func (sc *ServiceCreate) AddSubdomains(s ...*Subdomain) *ServiceCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sc.AddSubdomainIDs(ids...)
}

// Mutation returns the ServiceMutation object of the builder.
func (sc *ServiceCreate) Mutation() *ServiceMutation {
	return sc.mutation
}

// Save creates the Service in the database.
func (sc *ServiceCreate) Save(ctx context.Context) (*Service, error) {
	sc.defaults()
	return withHooks[*Service, ServiceMutation](ctx, sc.sqlSave, sc.mutation, sc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sc *ServiceCreate) SaveX(ctx context.Context) *Service {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *ServiceCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *ServiceCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sc *ServiceCreate) defaults() {
	if _, ok := sc.mutation.CreatedAt(); !ok {
		v := service.DefaultCreatedAt()
		sc.mutation.SetCreatedAt(v)
	}
	if _, ok := sc.mutation.UpdatedAt(); !ok {
		v := service.DefaultUpdatedAt()
		sc.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *ServiceCreate) check() error {
	if _, ok := sc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Service.name"`)}
	}
	if _, ok := sc.mutation.Protocol(); !ok {
		return &ValidationError{Name: "protocol", err: errors.New(`ent: missing required field "Service.protocol"`)}
	}
	if _, ok := sc.mutation.Domain(); !ok {
		return &ValidationError{Name: "domain", err: errors.New(`ent: missing required field "Service.domain"`)}
	}
	if _, ok := sc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Service.created_at"`)}
	}
	if _, ok := sc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Service.updated_at"`)}
	}
	return nil
}

func (sc *ServiceCreate) sqlSave(ctx context.Context) (*Service, error) {
	if err := sc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	sc.mutation.id = &_node.ID
	sc.mutation.done = true
	return _node, nil
}

func (sc *ServiceCreate) createSpec() (*Service, *sqlgraph.CreateSpec) {
	var (
		_node = &Service{config: sc.config}
		_spec = sqlgraph.NewCreateSpec(service.Table, sqlgraph.NewFieldSpec(service.FieldID, field.TypeInt))
	)
	if value, ok := sc.mutation.Name(); ok {
		_spec.SetField(service.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := sc.mutation.Protocol(); ok {
		_spec.SetField(service.FieldProtocol, field.TypeString, value)
		_node.Protocol = value
	}
	if value, ok := sc.mutation.Domain(); ok {
		_spec.SetField(service.FieldDomain, field.TypeString, value)
		_node.Domain = value
	}
	if value, ok := sc.mutation.CreatedAt(); ok {
		_spec.SetField(service.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := sc.mutation.UpdatedAt(); ok {
		_spec.SetField(service.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if nodes := sc.mutation.SubdomainsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.SubdomainsTable,
			Columns: []string{service.SubdomainsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: subdomain.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// ServiceCreateBulk is the builder for creating many Service entities in bulk.
type ServiceCreateBulk struct {
	config
	builders []*ServiceCreate
}

// Save creates the Service entities in the database.
func (scb *ServiceCreateBulk) Save(ctx context.Context) ([]*Service, error) {
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Service, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ServiceMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *ServiceCreateBulk) SaveX(ctx context.Context) []*Service {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *ServiceCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *ServiceCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}