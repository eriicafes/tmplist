package classic_pages

type FlashMessage struct {
	Message string
	Success bool
}

type Layout struct {
	Flash FlashMessage
}
