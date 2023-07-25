package db

import (
	"context"

	"hus-auth/ent"
	"hus-auth/ent/hussession"

	"github.com/google/uuid"
)

func QuerySessionBySID(c context.Context, client *ent.Client, sid string) (*ent.HusSession, error) {
	sid_uuid, err := uuid.Parse(sid)
	if err != nil {
		return nil, err
	}
	hs, err := client.HusSession.Query().Where(hussession.ID(sid_uuid)).WithConnectedSession().WithUser().Only(c)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}
	return hs, nil
}
