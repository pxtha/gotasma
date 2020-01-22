package types

const (
	PolicyObjectAny  = "*"
	PolicyActionAny  = "*"
	PolicySubjectAny = "*"
)

// Policy of user api
const (
	PolicyObjectUser         = "/api/v1/users/"
	PolicyActionUserReadList = "GET"

	PolicyObjectDev       = "/api/v1/users/dev"
	PolicyActionDevCreate = "POST"

	PolicyObjectDeleteDev = "/api/v1/users/dev"
	PolicyActionDevDelete = "DELETE"
)
