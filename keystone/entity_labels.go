package keystone

import (
	"github.com/kubex/keystone-go/proto"
)

type EntityLabelProvider interface {
	ClearKeystoneLabels() error
	GetKeystoneLabels() []*proto.EntityLabel
}

type EntityLabels struct {
	ksEntityLabels []*proto.EntityLabel
}

func (e *EntityLabels) ClearKeystoneLabels() error {
	e.ksEntityLabels = []*proto.EntityLabel{}
	return nil
}

func (e *EntityLabels) GetKeystoneLabels() []*proto.EntityLabel {
	return e.ksEntityLabels
}

func (e *EntityLabels) AddKeystoneLabel(name, value string) {
	e.ksEntityLabels = append(e.ksEntityLabels, &proto.EntityLabel{
		Name:  name,
		Value: value,
	})
}
