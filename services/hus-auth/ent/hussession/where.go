// Code generated by ent, DO NOT EDIT.

package hussession

import (
	"hus-auth/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldLTE(FieldID, id))
}

// Tid applies equality check predicate on the "tid" field. It's identical to TidEQ.
func Tid(v uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldTid, v))
}

// Iat applies equality check predicate on the "iat" field. It's identical to IatEQ.
func Iat(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldIat, v))
}

// Exp applies equality check predicate on the "exp" field. It's identical to ExpEQ.
func Exp(v int64) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldExp, v))
}

// Preserved applies equality check predicate on the "preserved" field. It's identical to PreservedEQ.
func Preserved(v bool) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldPreserved, v))
}

// UID applies equality check predicate on the "uid" field. It's identical to UIDEQ.
func UID(v uint64) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldUID, v))
}

// SignedAt applies equality check predicate on the "signed_at" field. It's identical to SignedAtEQ.
func SignedAt(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldSignedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldUpdatedAt, v))
}

// TidEQ applies the EQ predicate on the "tid" field.
func TidEQ(v uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldTid, v))
}

// TidNEQ applies the NEQ predicate on the "tid" field.
func TidNEQ(v uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldNEQ(FieldTid, v))
}

// TidIn applies the In predicate on the "tid" field.
func TidIn(vs ...uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldIn(FieldTid, vs...))
}

// TidNotIn applies the NotIn predicate on the "tid" field.
func TidNotIn(vs ...uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldNotIn(FieldTid, vs...))
}

// TidGT applies the GT predicate on the "tid" field.
func TidGT(v uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldGT(FieldTid, v))
}

// TidGTE applies the GTE predicate on the "tid" field.
func TidGTE(v uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldGTE(FieldTid, v))
}

// TidLT applies the LT predicate on the "tid" field.
func TidLT(v uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldLT(FieldTid, v))
}

// TidLTE applies the LTE predicate on the "tid" field.
func TidLTE(v uuid.UUID) predicate.HusSession {
	return predicate.HusSession(sql.FieldLTE(FieldTid, v))
}

// IatEQ applies the EQ predicate on the "iat" field.
func IatEQ(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldIat, v))
}

// IatNEQ applies the NEQ predicate on the "iat" field.
func IatNEQ(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldNEQ(FieldIat, v))
}

// IatIn applies the In predicate on the "iat" field.
func IatIn(vs ...time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldIn(FieldIat, vs...))
}

// IatNotIn applies the NotIn predicate on the "iat" field.
func IatNotIn(vs ...time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldNotIn(FieldIat, vs...))
}

// IatGT applies the GT predicate on the "iat" field.
func IatGT(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldGT(FieldIat, v))
}

// IatGTE applies the GTE predicate on the "iat" field.
func IatGTE(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldGTE(FieldIat, v))
}

// IatLT applies the LT predicate on the "iat" field.
func IatLT(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldLT(FieldIat, v))
}

// IatLTE applies the LTE predicate on the "iat" field.
func IatLTE(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldLTE(FieldIat, v))
}

// ExpEQ applies the EQ predicate on the "exp" field.
func ExpEQ(v int64) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldExp, v))
}

// ExpNEQ applies the NEQ predicate on the "exp" field.
func ExpNEQ(v int64) predicate.HusSession {
	return predicate.HusSession(sql.FieldNEQ(FieldExp, v))
}

// ExpIn applies the In predicate on the "exp" field.
func ExpIn(vs ...int64) predicate.HusSession {
	return predicate.HusSession(sql.FieldIn(FieldExp, vs...))
}

// ExpNotIn applies the NotIn predicate on the "exp" field.
func ExpNotIn(vs ...int64) predicate.HusSession {
	return predicate.HusSession(sql.FieldNotIn(FieldExp, vs...))
}

// ExpGT applies the GT predicate on the "exp" field.
func ExpGT(v int64) predicate.HusSession {
	return predicate.HusSession(sql.FieldGT(FieldExp, v))
}

// ExpGTE applies the GTE predicate on the "exp" field.
func ExpGTE(v int64) predicate.HusSession {
	return predicate.HusSession(sql.FieldGTE(FieldExp, v))
}

// ExpLT applies the LT predicate on the "exp" field.
func ExpLT(v int64) predicate.HusSession {
	return predicate.HusSession(sql.FieldLT(FieldExp, v))
}

// ExpLTE applies the LTE predicate on the "exp" field.
func ExpLTE(v int64) predicate.HusSession {
	return predicate.HusSession(sql.FieldLTE(FieldExp, v))
}

// PreservedEQ applies the EQ predicate on the "preserved" field.
func PreservedEQ(v bool) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldPreserved, v))
}

// PreservedNEQ applies the NEQ predicate on the "preserved" field.
func PreservedNEQ(v bool) predicate.HusSession {
	return predicate.HusSession(sql.FieldNEQ(FieldPreserved, v))
}

// UIDEQ applies the EQ predicate on the "uid" field.
func UIDEQ(v uint64) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldUID, v))
}

// UIDNEQ applies the NEQ predicate on the "uid" field.
func UIDNEQ(v uint64) predicate.HusSession {
	return predicate.HusSession(sql.FieldNEQ(FieldUID, v))
}

// UIDIn applies the In predicate on the "uid" field.
func UIDIn(vs ...uint64) predicate.HusSession {
	return predicate.HusSession(sql.FieldIn(FieldUID, vs...))
}

// UIDNotIn applies the NotIn predicate on the "uid" field.
func UIDNotIn(vs ...uint64) predicate.HusSession {
	return predicate.HusSession(sql.FieldNotIn(FieldUID, vs...))
}

// UIDIsNil applies the IsNil predicate on the "uid" field.
func UIDIsNil() predicate.HusSession {
	return predicate.HusSession(sql.FieldIsNull(FieldUID))
}

// UIDNotNil applies the NotNil predicate on the "uid" field.
func UIDNotNil() predicate.HusSession {
	return predicate.HusSession(sql.FieldNotNull(FieldUID))
}

// SignedAtEQ applies the EQ predicate on the "signed_at" field.
func SignedAtEQ(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldSignedAt, v))
}

// SignedAtNEQ applies the NEQ predicate on the "signed_at" field.
func SignedAtNEQ(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldNEQ(FieldSignedAt, v))
}

// SignedAtIn applies the In predicate on the "signed_at" field.
func SignedAtIn(vs ...time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldIn(FieldSignedAt, vs...))
}

// SignedAtNotIn applies the NotIn predicate on the "signed_at" field.
func SignedAtNotIn(vs ...time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldNotIn(FieldSignedAt, vs...))
}

// SignedAtGT applies the GT predicate on the "signed_at" field.
func SignedAtGT(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldGT(FieldSignedAt, v))
}

// SignedAtGTE applies the GTE predicate on the "signed_at" field.
func SignedAtGTE(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldGTE(FieldSignedAt, v))
}

// SignedAtLT applies the LT predicate on the "signed_at" field.
func SignedAtLT(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldLT(FieldSignedAt, v))
}

// SignedAtLTE applies the LTE predicate on the "signed_at" field.
func SignedAtLTE(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldLTE(FieldSignedAt, v))
}

// SignedAtIsNil applies the IsNil predicate on the "signed_at" field.
func SignedAtIsNil() predicate.HusSession {
	return predicate.HusSession(sql.FieldIsNull(FieldSignedAt))
}

// SignedAtNotNil applies the NotNil predicate on the "signed_at" field.
func SignedAtNotNil() predicate.HusSession {
	return predicate.HusSession(sql.FieldNotNull(FieldSignedAt))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.HusSession {
	return predicate.HusSession(sql.FieldLTE(FieldUpdatedAt, v))
}

// HasUser applies the HasEdge predicate on the "user" edge.
func HasUser() predicate.HusSession {
	return predicate.HusSession(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserWith applies the HasEdge predicate on the "user" edge with a given conditions (other predicates).
func HasUserWith(preds ...predicate.User) predicate.HusSession {
	return predicate.HusSession(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(UserInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasConnectedSession applies the HasEdge predicate on the "connected_session" edge.
func HasConnectedSession() predicate.HusSession {
	return predicate.HusSession(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ConnectedSessionTable, ConnectedSessionColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasConnectedSessionWith applies the HasEdge predicate on the "connected_session" edge with a given conditions (other predicates).
func HasConnectedSessionWith(preds ...predicate.ConnectedSession) predicate.HusSession {
	return predicate.HusSession(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ConnectedSessionInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ConnectedSessionTable, ConnectedSessionColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.HusSession) predicate.HusSession {
	return predicate.HusSession(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.HusSession) predicate.HusSession {
	return predicate.HusSession(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.HusSession) predicate.HusSession {
	return predicate.HusSession(func(s *sql.Selector) {
		p(s.Not())
	})
}
