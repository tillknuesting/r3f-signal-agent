package valueobject

type CollectorType string

const (
	CollectorTypeHTTP   CollectorType = "http"
	CollectorTypeRSS    CollectorType = "rss"
	CollectorTypeChrome CollectorType = "chrome"
)

func (c CollectorType) String() string {
	return string(c)
}
