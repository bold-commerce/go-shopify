package goshopify

type AccessScope struct {
	Handle string `json:"handle,omitempty"`
}

type AccessScopeResource struct {
	AccessScopes []AccessScope `json:"access_scopes"`
}

type AccessScopeService interface {
	List(options interface{}) ([]AccessScope, error)
}

type AccessScopeServiceOp struct {
	client *Client
}

func (a *AccessScopeServiceOp) List(options interface{}) ([]AccessScope, error) {
	resource := new(AccessScopeResource)
	err := a.client.Get("oauth/access_scopes.json", resource, nil)
	return resource.AccessScopes, err
}
