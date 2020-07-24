package kibela

type User struct {
	ID          `json:"id"`
	Account     string `json:"account"`
	AvatarImage struct {
		URL string `json:"url"`
	} `json:"avatarImage"`
}
