package helpers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SessionSet(c *gin.Context, UserID uint64) error {
	session := sessions.Default(c)
	session.Set("id", UserID)
	return session.Save()
}

func SessionGet(c *gin.Context) (uint64, bool) {
	session := sessions.Default(c)
	v := session.Get("id")
	if v == nil {
		return 0, false
	}
	id, ok := v.(uint64)
	return id, ok
}

func SessionClear(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()
	return session.Save()
}
