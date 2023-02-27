// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"hus-auth/ent/group"
	"hus-auth/ent/user"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// UserCreate is the builder for creating a User entity.
type UserCreate struct {
	config
	mutation *UserMutation
	hooks    []Hook
}

// SetUUID sets the "uuid" field.
func (uc *UserCreate) SetUUID(u uuid.UUID) *UserCreate {
	uc.mutation.SetUUID(u)
	return uc
}

// SetNillableUUID sets the "uuid" field if the given value is not nil.
func (uc *UserCreate) SetNillableUUID(u *uuid.UUID) *UserCreate {
	if u != nil {
		uc.SetUUID(*u)
	}
	return uc
}

// SetGoogleSub sets the "google_sub" field.
func (uc *UserCreate) SetGoogleSub(s string) *UserCreate {
	uc.mutation.SetGoogleSub(s)
	return uc
}

// SetEmail sets the "email" field.
func (uc *UserCreate) SetEmail(s string) *UserCreate {
	uc.mutation.SetEmail(s)
	return uc
}

// SetEmailVerified sets the "email_verified" field.
func (uc *UserCreate) SetEmailVerified(b bool) *UserCreate {
	uc.mutation.SetEmailVerified(b)
	return uc
}

// SetName sets the "name" field.
func (uc *UserCreate) SetName(s string) *UserCreate {
	uc.mutation.SetName(s)
	return uc
}

// SetBirthday sets the "birthday" field.
func (uc *UserCreate) SetBirthday(t time.Time) *UserCreate {
	uc.mutation.SetBirthday(t)
	return uc
}

// SetGivenName sets the "given_name" field.
func (uc *UserCreate) SetGivenName(s string) *UserCreate {
	uc.mutation.SetGivenName(s)
	return uc
}

// SetFamilyName sets the "family_name" field.
func (uc *UserCreate) SetFamilyName(s string) *UserCreate {
	uc.mutation.SetFamilyName(s)
	return uc
}

// SetGoogleProfilePicture sets the "google_profile_picture" field.
func (uc *UserCreate) SetGoogleProfilePicture(s string) *UserCreate {
	uc.mutation.SetGoogleProfilePicture(s)
	return uc
}

// AddGroupIDs adds the "groups" edge to the Group entity by IDs.
func (uc *UserCreate) AddGroupIDs(ids ...int) *UserCreate {
	uc.mutation.AddGroupIDs(ids...)
	return uc
}

// AddGroups adds the "groups" edges to the Group entity.
func (uc *UserCreate) AddGroups(g ...*Group) *UserCreate {
	ids := make([]int, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return uc.AddGroupIDs(ids...)
}

// Mutation returns the UserMutation object of the builder.
func (uc *UserCreate) Mutation() *UserMutation {
	return uc.mutation
}

// Save creates the User in the database.
func (uc *UserCreate) Save(ctx context.Context) (*User, error) {
	uc.defaults()
	return withHooks[*User, UserMutation](ctx, uc.sqlSave, uc.mutation, uc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (uc *UserCreate) SaveX(ctx context.Context) *User {
	v, err := uc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (uc *UserCreate) Exec(ctx context.Context) error {
	_, err := uc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uc *UserCreate) ExecX(ctx context.Context) {
	if err := uc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uc *UserCreate) defaults() {
	if _, ok := uc.mutation.UUID(); !ok {
		v := user.DefaultUUID()
		uc.mutation.SetUUID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uc *UserCreate) check() error {
	if _, ok := uc.mutation.UUID(); !ok {
		return &ValidationError{Name: "uuid", err: errors.New(`ent: missing required field "User.uuid"`)}
	}
	if _, ok := uc.mutation.GoogleSub(); !ok {
		return &ValidationError{Name: "google_sub", err: errors.New(`ent: missing required field "User.google_sub"`)}
	}
	if _, ok := uc.mutation.Email(); !ok {
		return &ValidationError{Name: "email", err: errors.New(`ent: missing required field "User.email"`)}
	}
	if _, ok := uc.mutation.EmailVerified(); !ok {
		return &ValidationError{Name: "email_verified", err: errors.New(`ent: missing required field "User.email_verified"`)}
	}
	if _, ok := uc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "User.name"`)}
	}
	if _, ok := uc.mutation.Birthday(); !ok {
		return &ValidationError{Name: "birthday", err: errors.New(`ent: missing required field "User.birthday"`)}
	}
	if _, ok := uc.mutation.GivenName(); !ok {
		return &ValidationError{Name: "given_name", err: errors.New(`ent: missing required field "User.given_name"`)}
	}
	if _, ok := uc.mutation.FamilyName(); !ok {
		return &ValidationError{Name: "family_name", err: errors.New(`ent: missing required field "User.family_name"`)}
	}
	if _, ok := uc.mutation.GoogleProfilePicture(); !ok {
		return &ValidationError{Name: "google_profile_picture", err: errors.New(`ent: missing required field "User.google_profile_picture"`)}
	}
	return nil
}

func (uc *UserCreate) sqlSave(ctx context.Context) (*User, error) {
	if err := uc.check(); err != nil {
		return nil, err
	}
	_node, _spec := uc.createSpec()
	if err := sqlgraph.CreateNode(ctx, uc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	uc.mutation.id = &_node.ID
	uc.mutation.done = true
	return _node, nil
}

func (uc *UserCreate) createSpec() (*User, *sqlgraph.CreateSpec) {
	var (
		_node = &User{config: uc.config}
		_spec = sqlgraph.NewCreateSpec(user.Table, sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt))
	)
	if value, ok := uc.mutation.UUID(); ok {
		_spec.SetField(user.FieldUUID, field.TypeUUID, value)
		_node.UUID = value
	}
	if value, ok := uc.mutation.GoogleSub(); ok {
		_spec.SetField(user.FieldGoogleSub, field.TypeString, value)
		_node.GoogleSub = value
	}
	if value, ok := uc.mutation.Email(); ok {
		_spec.SetField(user.FieldEmail, field.TypeString, value)
		_node.Email = value
	}
	if value, ok := uc.mutation.EmailVerified(); ok {
		_spec.SetField(user.FieldEmailVerified, field.TypeBool, value)
		_node.EmailVerified = value
	}
	if value, ok := uc.mutation.Name(); ok {
		_spec.SetField(user.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := uc.mutation.Birthday(); ok {
		_spec.SetField(user.FieldBirthday, field.TypeTime, value)
		_node.Birthday = value
	}
	if value, ok := uc.mutation.GivenName(); ok {
		_spec.SetField(user.FieldGivenName, field.TypeString, value)
		_node.GivenName = value
	}
	if value, ok := uc.mutation.FamilyName(); ok {
		_spec.SetField(user.FieldFamilyName, field.TypeString, value)
		_node.FamilyName = value
	}
	if value, ok := uc.mutation.GoogleProfilePicture(); ok {
		_spec.SetField(user.FieldGoogleProfilePicture, field.TypeString, value)
		_node.GoogleProfilePicture = value
	}
	if nodes := uc.mutation.GroupsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   user.GroupsTable,
			Columns: user.GroupsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: group.FieldID,
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

// UserCreateBulk is the builder for creating many User entities in bulk.
type UserCreateBulk struct {
	config
	builders []*UserCreate
}

// Save creates the User entities in the database.
func (ucb *UserCreateBulk) Save(ctx context.Context) ([]*User, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ucb.builders))
	nodes := make([]*User, len(ucb.builders))
	mutators := make([]Mutator, len(ucb.builders))
	for i := range ucb.builders {
		func(i int, root context.Context) {
			builder := ucb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*UserMutation)
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
					_, err = mutators[i+1].Mutate(root, ucb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ucb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ucb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ucb *UserCreateBulk) SaveX(ctx context.Context) []*User {
	v, err := ucb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ucb *UserCreateBulk) Exec(ctx context.Context) error {
	_, err := ucb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ucb *UserCreateBulk) ExecX(ctx context.Context) {
	if err := ucb.Exec(ctx); err != nil {
		panic(err)
	}
}