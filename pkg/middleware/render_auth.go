package middleware

import (
	"sync"

	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/util"
)

var renderKeysLock sync.Mutex
var renderKeys map[string]*m.SignedInUser = make(map[string]*m.SignedInUser)

func AddRenderAuthKey(orgId int64, userId int64, orgRole m.RoleType) (string, error) {
	renderKeysLock.Lock()
	defer renderKeysLock.Unlock()

	key, err := util.GetRandomString(32)
	if err != nil {
		return "", err
	}

	renderKeys[key] = &m.SignedInUser{
		OrgId:   orgId,
		OrgRole: orgRole,
		UserId:  userId,
	}

	return key, nil
}

func RemoveRenderAuthKey(key string) {
	renderKeysLock.Lock()
	defer renderKeysLock.Unlock()

	delete(renderKeys, key)
}
