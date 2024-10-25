package utils

type (
	// ResourceConfiguration is a struct for a configuration resource
	ResourceConfiguration struct {
		Path string
		Type ConfigResourceType
		New  func() Resource
	}

	Resource interface {
		GetID() string
	}

	// ConfigResourceType is the type of configuration resource
	ConfigResourceType string
)

const (
	ConfigResourceType_Object ConfigResourceType = "OBJECT"
	ConfigResourceType_Array  ConfigResourceType = "ARRAY"
)
