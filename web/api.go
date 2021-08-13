package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/internal/sapsystem"
	"github.com/trento-project/trento/internal/tags"
)

func ApiPingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

type JSONTag struct {
	Tag string `json:"tag" binding:"required"`
}

// ApiClusterCreateTagHandler godoc
// @Summary Add tag to Cluster
// @Accept json
// @Produce json
// @Param id path string true "Cluster id"
// @Param Body body JSONTag true "The tag to create"
// @Success 201 {object} JSONTag
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/clusters/{id}/tags [post]
func ApiClusterCreateTagHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		clusters, err := cluster.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if _, ok := clusters[id]; !ok {
			_ = c.Error(NotFoundError("could not find cluster"))
			return
		}

		var r JSONTag

		err = c.BindJSON(&r)
		if err != nil {
			_ = c.Error(BadRequestError("problems parsing JSON"))
			return
		}

		t := tags.NewTags(client, "clusters", id)
		err = t.Create(r.Tag)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, &r)
	}
}

// ApiClusterDeleteTagHandler godoc
// @Summary Delete a specific tag that belongs to a cluster
// @Accept json
// @Produce json
// @Param cluster path string true "Cluster id"
// @Param tag path string true "Tag"
// @Success 204 {object} map[string]interface{}
// @Router /api/clusters/{name}/tags/{tag} [delete]
func ApiClusterDeleteTagHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		tag := c.Param("tag")

		clusters, err := cluster.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if _, ok := clusters[id]; !ok {
			_ = c.Error(NotFoundError("could not find cluster"))
			return
		}

		t := tags.NewTags(client, "clusters", id)
		err = t.Delete(tag)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

// ApiSAPSystemCreateTagHandler godoc
// @Summary Add tag to SAPSystem
// @Accept json
// @Produce json
// @Param id path string true "SAPSystem id"
// @Param Body body JSONTag true "The tag to create"
// @Success 201 {object} JSONTag
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/sapsystems/{id}/tags [post]
func ApiSAPSystemCreateTagHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")

		// TODO: store sapsystem outside hosts
		hostList, err := hosts.Load(client, "", nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var system *sapsystem.SAPSystem
		for _, h := range hostList {
			sapSystems, err := h.GetSAPSystems()
			if err != nil {
				_ = c.Error(err)
				return
			}

			for _, s := range sapSystems {
				if s.SID == sid {
					system = s
					break
				}
			}
		}

		if system == nil {
			_ = c.Error(NotFoundError("could not find system"))
			return
		}

		var r JSONTag

		err = c.BindJSON(&r)
		if err != nil {
			_ = c.Error(BadRequestError("problems parsing JSON"))
			return
		}

		t := tags.NewTags(client, "sapsystems", sid)
		err = t.Create(r.Tag)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, &r)
	}
}

// ApiSAPSystemDeleteTagHandler godoc
// @Summary Delete a specific tag that belongs to a SAPSystem
// @Accept json
// @Produce json
// @Param cluster path string true "SAPSystem id"
// @Param tag path string true "Tag"
// @Success 204 {object} map[string]interface{}
// @Router /api/sapsystems/{name}/tags/{tag} [delete]
func ApiSAPSystemDeleteTagHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")
		tag := c.Param("tag")

		// TODO: store sapsystem outside hosts
		hostList, err := hosts.Load(client, "", nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var system *sapsystem.SAPSystem
		for _, h := range hostList {
			sapSystems, err := h.GetSAPSystems()
			if err != nil {
				_ = c.Error(err)
				return
			}

			for _, s := range sapSystems {
				if s.SID == sid {
					system = s
					break
				}
			}
		}

		if system == nil {
			_ = c.Error(NotFoundError("could not find system"))
			return
		}

		t := tags.NewTags(client, "sapsystems", sid)
		err = t.Delete(tag)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
