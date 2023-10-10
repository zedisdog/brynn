package ginx

type ContentType string

const (
	ContentTypeJson          ContentType = "application/json"
	ContentTypeForm          ContentType = "application/x-www-form-urlencoded"
	ContentTypeMultiPartForm ContentType = "multipart/form-data"
	ContentTypeXml           ContentType = "text/xml"
)
