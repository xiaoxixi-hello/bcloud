package controller

import (
	"gorm.io/gorm"
	"strings"
)

type Authority struct {
	ID    int `gorm:"primarykey"`
	Path  string
	Count int
}

func (a *App) GetAuthorityCount() int {
	var auth Authority
	err := a.CB.Debug().Take(&auth, "path = ?", strings.TrimLeft(a.Config["PanID"], "/")).Error
	if err == gorm.ErrRecordNotFound {
		a.CB.Create(&Authority{
			Path:  strings.TrimLeft(a.Config["PanID"], "/"),
			Count: 2,
		})
		return 2
	}
	return auth.Count
}
