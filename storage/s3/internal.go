package s3

type errStr string

const (
	ErrObjectTooLarge = errStr("object size exceeds maximum allowed by this interface")
	ErrBucketTooLarge = errStr("bucket size exceeds maximum allowed by this interface")
)

func (e errStr) Error() string {
	return string(e)
}
