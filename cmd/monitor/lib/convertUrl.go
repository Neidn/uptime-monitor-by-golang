package lib

import (
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"github.com/gosimple/slug"
)

func GetSlug(site config.Site) string {
	if site.Slug != "" {
		return site.Slug
	}
	return slug.Make(site.Name)
}
