package keystone

import (
	"github.com/kubex/keystone-go/proto"
)

// EntityLinkProvider is an interface for entities that can have links
type EntityLinkProvider interface {
	ClearKeystoneLinks() error
	GetKeystoneLinks() []*proto.EntityLink
	SetKeystoneLinks(links []*proto.EntityLink)
}

// EntityLinks is a struct that implements EntityLinkProvider
type EntityLinks struct {
	ksEntityLinks []*proto.EntityLink
}

// ClearKeystoneLinks clears the links
func (e *EntityLinks) ClearKeystoneLinks() error {
	e.ksEntityLinks = []*proto.EntityLink{}
	return nil
}

// GetKeystoneLinks returns the links
func (e *EntityLinks) GetKeystoneLinks() []*proto.EntityLink {
	return e.ksEntityLinks
}

// SetKeystoneLinks sets the links
func (e *EntityLinks) SetKeystoneLinks(links []*proto.EntityLink) {
	e.ksEntityLinks = links
}

// AddKeystoneLink adds a link
func (e *EntityLinks) AddKeystoneLink(key, name, location string) {
	e.ksEntityLinks = append(e.ksEntityLinks, &proto.EntityLink{
		Type:     &proto.Key{Key: key},
		Name:     name,
		Location: location,
	})
}
