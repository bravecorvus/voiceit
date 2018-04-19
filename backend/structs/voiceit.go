// File voiceit.go defines the Go structs in which to store returned JSON objects from the api.voiceit endpoints
package structs

type CreateGroupResponse struct {
	GroupID string `json:"groupId"`
}

type CreateNewUserResponse struct {
	UserID       string `json:"userId"`
	ResponseCode string `json:"resonseCode"`
}

type AddUserToGroupResponse struct {
	ResponseCode string `json:"resonseCode"`
}

type GetAllGroupsResponse struct {
	ResponseCode string  `json:"responseCode"`
	Groups       []Group `json:"groups"`
}

type Group struct {
	Description string   `json:"description"`
	GroupID     string   `json:"groupId"`
	Users       []string `json:"users"`
}

type GetSpecificGroupResponse struct {
	ResponseCode string   `json:"responseCode"`
	Users        []string `json:"users"`
}

type RemoveUserFromGroupResponse struct {
	ResponseCode string `json:"responseCode"`
}
