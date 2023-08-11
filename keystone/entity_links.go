package keystone

import (
	"github.com/kubex/keystone-go/proto"
)

type EntityLinkProvider interface {
	ClearKeystoneLinks() error
	GetKeystoneLinks() []*proto.EntityLink
}

type EntityLinks struct {
	ksEntityLinks []*proto.EntityLink
}

func (e *EntityLinks) ClearKeystoneLinks() error {
	e.ksEntityLinks = []*proto.EntityLink{}
	return nil
}

func (e *EntityLinks) GetKeystoneLinks() []*proto.EntityLink {
	return e.ksEntityLinks
}

func (e *EntityLinks) AddKeystoneLink(key, name, location string) {
	e.ksEntityLinks = append(e.ksEntityLinks, &proto.EntityLink{
		Type:     &proto.Key{Key: key},
		Name:     name,
		Location: location,
	})
}
