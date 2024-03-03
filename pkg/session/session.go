package session

import (
	"errors"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	cache "github.com/patrickmn/go-cache"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/sessionid"
)

// nonceおよびstateは、OIDCのフローの実行中だけ保持しておけば良いため、短い
const temporaryExpireTime time.Duration = 3 * time.Minute

// IDTokenおよびUserInfoは、ログインしている間は保持し続ける必要があるため、長い
const sessionExpireTime time.Duration = 1 * time.Hour

// 5分でExpireするデフォルト設定は、シグネチャが要求するため設定しているが、実際は使っていない
// 手動でキャッシュ時間を設定している
var dataStore *cache.Cache

func Init() {
	dataStore = cache.New(5*time.Minute, 5*time.Minute)
}

func getNonceKey(id sessionid.ID) string {
	return string(id + "nonce")
}
func getStateKey(id sessionid.ID) string {
	return string(id + "state")
}
func getRedirectURLKey(id sessionid.ID) string {
	return string(id + "redirectURL")
}
func getIDTokenKey(id sessionid.ID) string {
	return string(id + "IDToken")
}
func getUserInfoKey(id sessionid.ID) string {
	return string(id + "UserInfo")
}

func SetNonce(id sessionid.ID, nonce string) error {
	key := getNonceKey(id)
	return dataStore.Add(key, nonce, temporaryExpireTime)
}

func SetState(id sessionid.ID, state string) error {
	key := getStateKey(id)
	return dataStore.Add(key, state, temporaryExpireTime)
}

func SetRedirectURL(id sessionid.ID, redirectURL string) error {
	key := getRedirectURLKey(id)
	return dataStore.Add(key, redirectURL, temporaryExpireTime)
}

func SetIDToken(id sessionid.ID, idToken *oidc.IDToken) error {
	key := getIDTokenKey(id)
	return dataStore.Add(key, idToken, sessionExpireTime)
}

func SetUserInfo(id sessionid.ID, userInfo *oidc.UserInfo) error {
	key := getUserInfoKey(id)
	return dataStore.Add(key, userInfo, sessionExpireTime)
}

func DeleteNonce(id sessionid.ID) {
	key := getNonceKey(id)
	dataStore.Delete(key)
}

func DeleteState(id sessionid.ID) {
	key := getNonceKey(id)
	dataStore.Delete(key)
}

func DeleteRedirectURL(id sessionid.ID) {
	key := getRedirectURLKey(id)
	dataStore.Delete(key)
}

func DeleteIDToken(id sessionid.ID) {
	key := getIDTokenKey(id)
	dataStore.Delete(key)
}

func DeleteUserInfo(id sessionid.ID) {
	key := getUserInfoKey(id)
	dataStore.Delete(key)
}

func GetNonce(id sessionid.ID) (string, error) {
	key := getNonceKey(id)
	nonce, found := dataStore.Get(key)
	if !found {
		return "", errors.New("error: nonce not found")
	}
	return nonce.(string), nil
}

func GetState(id sessionid.ID) (string, error) {
	key := getStateKey(id)
	state, found := dataStore.Get(key)
	if !found {
		return "", errors.New("error: state not found")
	}
	return state.(string), nil
}

func GetRedirectURL(id sessionid.ID) (string, error) {
	key := getRedirectURLKey(id)
	redirectURL, found := dataStore.Get(key)
	if !found {
		return "", errors.New("error: redirect URL not found")
	}
	return redirectURL.(string), nil
}

func GetIDToken(id sessionid.ID) (*oidc.IDToken, error) {
	key := getIDTokenKey(id)
	idToken, found := dataStore.Get(key)
	if !found {
		return nil, errors.New("error: IDToken not found")
	}
	return idToken.(*oidc.IDToken), nil
}

func GetUserInfo(id sessionid.ID) (*oidc.UserInfo, error) {
	key := getUserInfoKey(id)
	userInfo, found := dataStore.Get(key)
	if !found {
		return nil, errors.New("error: UserInfo not found")
	}
	return userInfo.(*oidc.UserInfo), nil
}

func RefreshSession(oldID, newID sessionid.ID) error {
	defer func() {
		// リフレッシュに失敗するような異常な事態では、最悪を避けるために安全側に倒す
		DeleteIDToken(oldID)
		DeleteUserInfo(oldID)
	}()

	if idToken, err := GetIDToken(oldID); err == nil {
		if setErr := SetIDToken(newID, idToken); setErr != nil {
			return setErr
		}
	}

	if userInfo, err := GetUserInfo(oldID); err == nil {
		if setErr := SetUserInfo(newID, userInfo); setErr != nil {
			return setErr
		}
	}

	return nil
}
