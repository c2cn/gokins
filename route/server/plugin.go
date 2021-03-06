package server

import (
	"fmt"
	"gokins/comm"
	"gokins/mgr"
	"gokins/model"
	"gokins/models"
	"gokins/service/dbService"

	"github.com/gin-gonic/gin"
	ruisUtil "github.com/mgr9525/go-ruisutil"
)

func PlugList(c *gin.Context, req *ruisUtil.Map) {
	tid, err := req.GetInt("tid")
	if err != nil || tid <= 0 {
		c.String(500, "param err")
		return
	}
	ls := make([]*model.TPlugin, 0)
	ses := comm.Db.Where("del!='1' and tid=?", tid).OrderBy("sort ASC,id ASC")
	err = ses.Find(&ls)
	if err != nil {
		c.String(500, "find err:"+err.Error())
		return
	}
	c.JSON(200, ls)
}
func PlugEdit(c *gin.Context, req *models.Plugin) {
	if req.Title == "" || req.Tid <= 0 {
		c.String(500, "param err")
		return
	}
	if err := req.Save(); err != nil {
		c.String(500, "save err:"+err.Error())
		return
	}
	c.String(200, fmt.Sprintf("%d", req.Id))
}
func PlugDel(c *gin.Context, req *ruisUtil.Map) {
	id, err := req.GetInt("id")
	if err != nil || id <= 0 {
		c.String(500, "param err")
		return
	}
	m := &models.Plugin{}
	if err := m.Del(int(id)); err != nil {
		c.String(500, "save err:"+err.Error())
		return
	}
	c.String(200, fmt.Sprintf("%d", m.Id))
}
func PlugRuns(c *gin.Context, req *ruisUtil.Map) {
	id, err := req.GetInt("id")
	if err != nil || id <= 0 {
		c.String(500, "param err")
		return
	}

	mr := dbService.GetModelRun(int(id))
	if mr == nil {
		c.String(404, "not found")
		return
	}

	ls := make([]*models.PluginUI, 0)
	ses := comm.Db.Where("del!='1' and tid=?", mr.Tid).OrderBy("sort ASC,id ASC")
	err = ses.Find(&ls)
	if err != nil {
		c.String(500, "find err:"+err.Error())
		return
	}
	for _, v := range ls {
		v.ToUI(mr.Id)
	}

	res := ruisUtil.NewMap()
	res.Set("list", ls)
	res.Set("tid", mr.Tid)
	res.Set("end", mr.State >= 2)
	c.JSON(200, res)
}
func PlugLog(c *gin.Context, req *ruisUtil.Map) {
	tid, err := req.GetInt("tid")
	if err != nil || tid <= 0 {
		c.String(500, "param err")
		return
	}
	pid, err := req.GetInt("pid")
	if err != nil || pid <= 0 {
		c.String(500, "param err")
		return
	}
	mr := dbService.GetModelRun(int(tid))
	if mr == nil {
		c.String(404, "not found")
		return
	}
	e := dbService.FindPluginRun(mr.Tid, mr.Id, int(pid))
	if e == nil {
		c.String(404, "not found")
		return
	}
	res := ruisUtil.NewMap()
	res.Set("up", true)
	res.Set("text", mgr.ExecMgr.TaskRead(mr.Id, e.Id))
	if e != nil && e.State >= 2 {
		res.Set("up", false)
		res.Set("text", e.Output)
	}
	c.JSON(200, res)
}
