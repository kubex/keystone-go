package keystone

import (
	"github.com/kubex/keystone-go/proto"
)

// EntityLabelProvider is an interface for entities that can provide labels
type EntityLabelProvider interface {
	ClearKeystoneLabels() error
	GetKeystoneLabels() []*proto.EntityLabel
}

// EntityLabels is a struct that implements EntityLabelProvider
type EntityLabels struct {
	ksEntityLabels []*proto.EntityLabel
}

// ClearKeystoneLabels clears the labels
func (e *EntityLabels) ClearKeystoneLabels() error {
	e.ksEntityLabels = []*proto.EntityLabel{}
	return nil
}

// GetKeystoneLabels returns the labels
func (e *EntityLabels) GetKeystoneLabels() []*proto.EntityLabel {
	return e.ksEntityLabels
}

// AddKeystoneLabel adds a label
func (e *EntityLabels) AddKeystoneLabel(name, value string) {
	e.ksEntityLabels = append(e.ksEntityLabels, &proto.EntityLabel{
		Name:  name,
		Value: value,
	})
}
