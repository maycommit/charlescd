package v1alpha1

type objectMeta struct {
	Name *string
}

func (a *Circle) GetMetadata() *objectMeta {
	var om objectMeta
	if a != nil {
		om.Name = &a.Name
	}
	return &om
}
