// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"hus-auth/ent/hussession"
	"hus-auth/ent/predicate"
	"hus-auth/ent/user"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// UserUpdate is the builder for updating User entities.
type UserUpdate struct {
	config
	hooks    []Hook
	mutation *UserMutation
}

// Where appends a list predicates to the UserUpdate builder.
func (uu *UserUpdate) Where(ps ...predicate.User) *UserUpdate {
	uu.mutation.Where(ps...)
	return uu
}

// SetProvider sets the "provider" field.
func (uu *UserUpdate) SetProvider(u user.Provider) *UserUpdate {
	uu.mutation.SetProvider(u)
	return uu
}

// SetGoogleSub sets the "google_sub" field.
func (uu *UserUpdate) SetGoogleSub(s string) *UserUpdate {
	uu.mutation.SetGoogleSub(s)
	return uu
}

// SetEmail sets the "email" field.
func (uu *UserUpdate) SetEmail(s string) *UserUpdate {
	uu.mutation.SetEmail(s)
	return uu
}

// SetEmailVerified sets the "email_verified" field.
func (uu *UserUpdate) SetEmailVerified(b bool) *UserUpdate {
	uu.mutation.SetEmailVerified(b)
	return uu
}

// SetName sets the "name" field.
func (uu *UserUpdate) SetName(s string) *UserUpdate {
	uu.mutation.SetName(s)
	return uu
}

// SetGivenName sets the "given_name" field.
func (uu *UserUpdate) SetGivenName(s string) *UserUpdate {
	uu.mutation.SetGivenName(s)
	return uu
}

// SetFamilyName sets the "family_name" field.
func (uu *UserUpdate) SetFamilyName(s string) *UserUpdate {
	uu.mutation.SetFamilyName(s)
	return uu
}

// SetBirthdate sets the "birthdate" field.
func (uu *UserUpdate) SetBirthdate(t time.Time) *UserUpdate {
	uu.mutation.SetBirthdate(t)
	return uu
}

// SetNillableBirthdate sets the "birthdate" field if the given value is not nil.
func (uu *UserUpdate) SetNillableBirthdate(t *time.Time) *UserUpdate {
	if t != nil {
		uu.SetBirthdate(*t)
	}
	return uu
}

// ClearBirthdate clears the value of the "birthdate" field.
func (uu *UserUpdate) ClearBirthdate() *UserUpdate {
	uu.mutation.ClearBirthdate()
	return uu
}

// SetProfilePictureURL sets the "profile_picture_url" field.
func (uu *UserUpdate) SetProfilePictureURL(s string) *UserUpdate {
	uu.mutation.SetProfilePictureURL(s)
	return uu
}

// SetNillableProfilePictureURL sets the "profile_picture_url" field if the given value is not nil.
func (uu *UserUpdate) SetNillableProfilePictureURL(s *string) *UserUpdate {
	if s != nil {
		uu.SetProfilePictureURL(*s)
	}
	return uu
}

// ClearProfilePictureURL clears the value of the "profile_picture_url" field.
func (uu *UserUpdate) ClearProfilePictureURL() *UserUpdate {
	uu.mutation.ClearProfilePictureURL()
	return uu
}

// SetCreatedAt sets the "created_at" field.
func (uu *UserUpdate) SetCreatedAt(t time.Time) *UserUpdate {
	uu.mutation.SetCreatedAt(t)
	return uu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (uu *UserUpdate) SetNillableCreatedAt(t *time.Time) *UserUpdate {
	if t != nil {
		uu.SetCreatedAt(*t)
	}
	return uu
}

// SetUpdatedAt sets the "updated_at" field.
func (uu *UserUpdate) SetUpdatedAt(t time.Time) *UserUpdate {
	uu.mutation.SetUpdatedAt(t)
	return uu
}

// AddHusSessionIDs adds the "hus_sessions" edge to the HusSession entity by IDs.
func (uu *UserUpdate) AddHusSessionIDs(ids ...uuid.UUID) *UserUpdate {
	uu.mutation.AddHusSessionIDs(ids...)
	return uu
}

// AddHusSessions adds the "hus_sessions" edges to the HusSession entity.
func (uu *UserUpdate) AddHusSessions(h ...*HusSession) *UserUpdate {
	ids := make([]uuid.UUID, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return uu.AddHusSessionIDs(ids...)
}

// Mutation returns the UserMutation object of the builder.
func (uu *UserUpdate) Mutation() *UserMutation {
	return uu.mutation
}

// ClearHusSessions clears all "hus_sessions" edges to the HusSession entity.
func (uu *UserUpdate) ClearHusSessions() *UserUpdate {
	uu.mutation.ClearHusSessions()
	return uu
}

// RemoveHusSessionIDs removes the "hus_sessions" edge to HusSession entities by IDs.
func (uu *UserUpdate) RemoveHusSessionIDs(ids ...uuid.UUID) *UserUpdate {
	uu.mutation.RemoveHusSessionIDs(ids...)
	return uu
}

// RemoveHusSessions removes "hus_sessions" edges to HusSession entities.
func (uu *UserUpdate) RemoveHusSessions(h ...*HusSession) *UserUpdate {
	ids := make([]uuid.UUID, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return uu.RemoveHusSessionIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (uu *UserUpdate) Save(ctx context.Context) (int, error) {
	uu.defaults()
	return withHooks[int, UserMutation](ctx, uu.sqlSave, uu.mutation, uu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (uu *UserUpdate) SaveX(ctx context.Context) int {
	affected, err := uu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (uu *UserUpdate) Exec(ctx context.Context) error {
	_, err := uu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uu *UserUpdate) ExecX(ctx context.Context) {
	if err := uu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uu *UserUpdate) defaults() {
	if _, ok := uu.mutation.UpdatedAt(); !ok {
		v := user.UpdateDefaultUpdatedAt()
		uu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uu *UserUpdate) check() error {
	if v, ok := uu.mutation.Provider(); ok {
		if err := user.ProviderValidator(v); err != nil {
			return &ValidationError{Name: "provider", err: fmt.Errorf(`ent: validator failed for field "User.provider": %w`, err)}
		}
	}
	return nil
}

func (uu *UserUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := uu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(user.Table, user.Columns, sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID))
	if ps := uu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uu.mutation.Provider(); ok {
		_spec.SetField(user.FieldProvider, field.TypeEnum, value)
	}
	if value, ok := uu.mutation.GoogleSub(); ok {
		_spec.SetField(user.FieldGoogleSub, field.TypeString, value)
	}
	if value, ok := uu.mutation.Email(); ok {
		_spec.SetField(user.FieldEmail, field.TypeString, value)
	}
	if value, ok := uu.mutation.EmailVerified(); ok {
		_spec.SetField(user.FieldEmailVerified, field.TypeBool, value)
	}
	if value, ok := uu.mutation.Name(); ok {
		_spec.SetField(user.FieldName, field.TypeString, value)
	}
	if value, ok := uu.mutation.GivenName(); ok {
		_spec.SetField(user.FieldGivenName, field.TypeString, value)
	}
	if value, ok := uu.mutation.FamilyName(); ok {
		_spec.SetField(user.FieldFamilyName, field.TypeString, value)
	}
	if value, ok := uu.mutation.Birthdate(); ok {
		_spec.SetField(user.FieldBirthdate, field.TypeTime, value)
	}
	if uu.mutation.BirthdateCleared() {
		_spec.ClearField(user.FieldBirthdate, field.TypeTime)
	}
	if value, ok := uu.mutation.ProfilePictureURL(); ok {
		_spec.SetField(user.FieldProfilePictureURL, field.TypeString, value)
	}
	if uu.mutation.ProfilePictureURLCleared() {
		_spec.ClearField(user.FieldProfilePictureURL, field.TypeString)
	}
	if value, ok := uu.mutation.CreatedAt(); ok {
		_spec.SetField(user.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := uu.mutation.UpdatedAt(); ok {
		_spec.SetField(user.FieldUpdatedAt, field.TypeTime, value)
	}
	if uu.mutation.HusSessionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.HusSessionsTable,
			Columns: []string{user.HusSessionsColumn},
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
	if nodes := uu.mutation.RemovedHusSessionsIDs(); len(nodes) > 0 && !uu.mutation.HusSessionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.HusSessionsTable,
			Columns: []string{user.HusSessionsColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.HusSessionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.HusSessionsTable,
			Columns: []string{user.HusSessionsColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, uu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	uu.mutation.done = true
	return n, nil
}

// UserUpdateOne is the builder for updating a single User entity.
type UserUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *UserMutation
}

// SetProvider sets the "provider" field.
func (uuo *UserUpdateOne) SetProvider(u user.Provider) *UserUpdateOne {
	uuo.mutation.SetProvider(u)
	return uuo
}

// SetGoogleSub sets the "google_sub" field.
func (uuo *UserUpdateOne) SetGoogleSub(s string) *UserUpdateOne {
	uuo.mutation.SetGoogleSub(s)
	return uuo
}

// SetEmail sets the "email" field.
func (uuo *UserUpdateOne) SetEmail(s string) *UserUpdateOne {
	uuo.mutation.SetEmail(s)
	return uuo
}

// SetEmailVerified sets the "email_verified" field.
func (uuo *UserUpdateOne) SetEmailVerified(b bool) *UserUpdateOne {
	uuo.mutation.SetEmailVerified(b)
	return uuo
}

// SetName sets the "name" field.
func (uuo *UserUpdateOne) SetName(s string) *UserUpdateOne {
	uuo.mutation.SetName(s)
	return uuo
}

// SetGivenName sets the "given_name" field.
func (uuo *UserUpdateOne) SetGivenName(s string) *UserUpdateOne {
	uuo.mutation.SetGivenName(s)
	return uuo
}

// SetFamilyName sets the "family_name" field.
func (uuo *UserUpdateOne) SetFamilyName(s string) *UserUpdateOne {
	uuo.mutation.SetFamilyName(s)
	return uuo
}

// SetBirthdate sets the "birthdate" field.
func (uuo *UserUpdateOne) SetBirthdate(t time.Time) *UserUpdateOne {
	uuo.mutation.SetBirthdate(t)
	return uuo
}

// SetNillableBirthdate sets the "birthdate" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableBirthdate(t *time.Time) *UserUpdateOne {
	if t != nil {
		uuo.SetBirthdate(*t)
	}
	return uuo
}

// ClearBirthdate clears the value of the "birthdate" field.
func (uuo *UserUpdateOne) ClearBirthdate() *UserUpdateOne {
	uuo.mutation.ClearBirthdate()
	return uuo
}

// SetProfilePictureURL sets the "profile_picture_url" field.
func (uuo *UserUpdateOne) SetProfilePictureURL(s string) *UserUpdateOne {
	uuo.mutation.SetProfilePictureURL(s)
	return uuo
}

// SetNillableProfilePictureURL sets the "profile_picture_url" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableProfilePictureURL(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetProfilePictureURL(*s)
	}
	return uuo
}

// ClearProfilePictureURL clears the value of the "profile_picture_url" field.
func (uuo *UserUpdateOne) ClearProfilePictureURL() *UserUpdateOne {
	uuo.mutation.ClearProfilePictureURL()
	return uuo
}

// SetCreatedAt sets the "created_at" field.
func (uuo *UserUpdateOne) SetCreatedAt(t time.Time) *UserUpdateOne {
	uuo.mutation.SetCreatedAt(t)
	return uuo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableCreatedAt(t *time.Time) *UserUpdateOne {
	if t != nil {
		uuo.SetCreatedAt(*t)
	}
	return uuo
}

// SetUpdatedAt sets the "updated_at" field.
func (uuo *UserUpdateOne) SetUpdatedAt(t time.Time) *UserUpdateOne {
	uuo.mutation.SetUpdatedAt(t)
	return uuo
}

// AddHusSessionIDs adds the "hus_sessions" edge to the HusSession entity by IDs.
func (uuo *UserUpdateOne) AddHusSessionIDs(ids ...uuid.UUID) *UserUpdateOne {
	uuo.mutation.AddHusSessionIDs(ids...)
	return uuo
}

// AddHusSessions adds the "hus_sessions" edges to the HusSession entity.
func (uuo *UserUpdateOne) AddHusSessions(h ...*HusSession) *UserUpdateOne {
	ids := make([]uuid.UUID, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return uuo.AddHusSessionIDs(ids...)
}

// Mutation returns the UserMutation object of the builder.
func (uuo *UserUpdateOne) Mutation() *UserMutation {
	return uuo.mutation
}

// ClearHusSessions clears all "hus_sessions" edges to the HusSession entity.
func (uuo *UserUpdateOne) ClearHusSessions() *UserUpdateOne {
	uuo.mutation.ClearHusSessions()
	return uuo
}

// RemoveHusSessionIDs removes the "hus_sessions" edge to HusSession entities by IDs.
func (uuo *UserUpdateOne) RemoveHusSessionIDs(ids ...uuid.UUID) *UserUpdateOne {
	uuo.mutation.RemoveHusSessionIDs(ids...)
	return uuo
}

// RemoveHusSessions removes "hus_sessions" edges to HusSession entities.
func (uuo *UserUpdateOne) RemoveHusSessions(h ...*HusSession) *UserUpdateOne {
	ids := make([]uuid.UUID, len(h))
	for i := range h {
		ids[i] = h[i].ID
	}
	return uuo.RemoveHusSessionIDs(ids...)
}

// Where appends a list predicates to the UserUpdate builder.
func (uuo *UserUpdateOne) Where(ps ...predicate.User) *UserUpdateOne {
	uuo.mutation.Where(ps...)
	return uuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (uuo *UserUpdateOne) Select(field string, fields ...string) *UserUpdateOne {
	uuo.fields = append([]string{field}, fields...)
	return uuo
}

// Save executes the query and returns the updated User entity.
func (uuo *UserUpdateOne) Save(ctx context.Context) (*User, error) {
	uuo.defaults()
	return withHooks[*User, UserMutation](ctx, uuo.sqlSave, uuo.mutation, uuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (uuo *UserUpdateOne) SaveX(ctx context.Context) *User {
	node, err := uuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (uuo *UserUpdateOne) Exec(ctx context.Context) error {
	_, err := uuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uuo *UserUpdateOne) ExecX(ctx context.Context) {
	if err := uuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uuo *UserUpdateOne) defaults() {
	if _, ok := uuo.mutation.UpdatedAt(); !ok {
		v := user.UpdateDefaultUpdatedAt()
		uuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uuo *UserUpdateOne) check() error {
	if v, ok := uuo.mutation.Provider(); ok {
		if err := user.ProviderValidator(v); err != nil {
			return &ValidationError{Name: "provider", err: fmt.Errorf(`ent: validator failed for field "User.provider": %w`, err)}
		}
	}
	return nil
}

func (uuo *UserUpdateOne) sqlSave(ctx context.Context) (_node *User, err error) {
	if err := uuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(user.Table, user.Columns, sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID))
	id, ok := uuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "User.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := uuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, user.FieldID)
		for _, f := range fields {
			if !user.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != user.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := uuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uuo.mutation.Provider(); ok {
		_spec.SetField(user.FieldProvider, field.TypeEnum, value)
	}
	if value, ok := uuo.mutation.GoogleSub(); ok {
		_spec.SetField(user.FieldGoogleSub, field.TypeString, value)
	}
	if value, ok := uuo.mutation.Email(); ok {
		_spec.SetField(user.FieldEmail, field.TypeString, value)
	}
	if value, ok := uuo.mutation.EmailVerified(); ok {
		_spec.SetField(user.FieldEmailVerified, field.TypeBool, value)
	}
	if value, ok := uuo.mutation.Name(); ok {
		_spec.SetField(user.FieldName, field.TypeString, value)
	}
	if value, ok := uuo.mutation.GivenName(); ok {
		_spec.SetField(user.FieldGivenName, field.TypeString, value)
	}
	if value, ok := uuo.mutation.FamilyName(); ok {
		_spec.SetField(user.FieldFamilyName, field.TypeString, value)
	}
	if value, ok := uuo.mutation.Birthdate(); ok {
		_spec.SetField(user.FieldBirthdate, field.TypeTime, value)
	}
	if uuo.mutation.BirthdateCleared() {
		_spec.ClearField(user.FieldBirthdate, field.TypeTime)
	}
	if value, ok := uuo.mutation.ProfilePictureURL(); ok {
		_spec.SetField(user.FieldProfilePictureURL, field.TypeString, value)
	}
	if uuo.mutation.ProfilePictureURLCleared() {
		_spec.ClearField(user.FieldProfilePictureURL, field.TypeString)
	}
	if value, ok := uuo.mutation.CreatedAt(); ok {
		_spec.SetField(user.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := uuo.mutation.UpdatedAt(); ok {
		_spec.SetField(user.FieldUpdatedAt, field.TypeTime, value)
	}
	if uuo.mutation.HusSessionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.HusSessionsTable,
			Columns: []string{user.HusSessionsColumn},
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
	if nodes := uuo.mutation.RemovedHusSessionsIDs(); len(nodes) > 0 && !uuo.mutation.HusSessionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.HusSessionsTable,
			Columns: []string{user.HusSessionsColumn},
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.HusSessionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.HusSessionsTable,
			Columns: []string{user.HusSessionsColumn},
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
	_node = &User{config: uuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, uuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	uuo.mutation.done = true
	return _node, nil
}