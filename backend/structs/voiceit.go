// File voiceit.go defines the Go structs in which to store returned JSON objects from the api.voiceit endpoints
package structs

type CreateNewUserResponse struct {
	UserID       string `json:"userId"`
	ResponseCode string `json:"responseCode"`
	Message      string `json:"message"`
}

type RemoveUserFromGroupResponse struct {
	ResponseCode string `json:"responseCode"`
}

type CreateUserVideoEnrollmentResponse struct {
	ResponseCode string `json:"responseCode"`
	Message      string `json:"message"`
}

type VideoVerificationResponse struct {
	ResponseCode string `json:"responseCode"`
	Message      string `json:"message"`
}
