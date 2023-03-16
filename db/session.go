package db

import (
	"context"
	"fmt"

	"hus-auth/ent"
	"hus-auth/ent/hussession"

	"github.com/google/uuid"
)

func QuerySessionBySID(c context.Context, client *ent.Client, sid string) (*ent.HusSession, error) {
	sid_uuid, err := uuid.Parse(sid)
	if err != nil {
		err = fmt.Errorf("[F]invalid sid:%v", err)
		return nil, err
	}
	return client.HusSession.Query().Where(hussession.ID(sid_uuid)).Only(c)
}
