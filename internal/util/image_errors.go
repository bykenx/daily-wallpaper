package util

type DownloadError struct {
	msg string
	error
}

func (e DownloadError) Error() string {
	return e.msg
}
