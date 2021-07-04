package agin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type Admin struct {
	JwtLive   time.Duration
	JwtSecret string `mapstructure:"jwtSecret" json:"jwtSecret" yarm:"jwtSecret"`

	OAuthConnPool ConnectionPool

	upgrade websocket.Upgrader
}

type register struct {
	Models map[string]IAdminModel
}

func NewRegister() register {
	return register{make(map[string]IAdminModel)}
}

type IAdminModel interface {
	GetColumnInfo() []ColumnInfoDetail
	GetRecord(c *gin.Context) (gin.H, error)
	GetRecordAmount() int64
	Updates(c *gin.Context)
	Create(c *gin.Context)
	Delete(c *gin.Context)
}

func (r *register) AddAdminModel(model IAdminModel, name string) {
	r.Models[name] = model
}

type ColumnInfoHandler struct {}


type ColumnInfoDetail struct {
	JsonName string						`json:"jsonName"`
	Comment string						`json:"comment"`
	Other OtherColumnInfo				`json:"other"`
}

type OtherColumnInfo struct {
	Kind string							`json:"kind"`
	Type string							`json:"type"`
	TypeInfo string						`json:"typeInfo"`
	More []string						`json:"more"`
}

func (c *OtherColumnInfo) addToMore(info string) {
	c.More = append(c.More, info)
}

func (c *OtherColumnInfo) GetInfo(adminTag string, gormTab string, kind reflect.Kind) {
	if kind == reflect.String {
		c.Kind = "string"
	} else {
		c.Kind = "number"
	}
	adminTypeRe, _ := regexp.Compile(`type:(.+[^:]):([^;,]+)`)
	//re_sub, _ := regexp.Compile(`sub:(.+[^;,])`)
	match := adminTypeRe.FindStringSubmatch(adminTag)
	if len(match) > 0 {
		c.Type = match[1]
		c.addToMore(match[2])

	} else {
		adminTypeRe, _ = regexp.Compile(`type:([^;,]+)`)
		match = adminTypeRe.FindStringSubmatch(adminTag)
		if len(match) > 0 {
			c.Type = match[1]
		} else {
			gormTypeRe, _ := regexp.Compile(`ENUM\(([^;]+)\)`)
			match = gormTypeRe.FindStringSubmatch(gormTab)
			if len(match) > 0 {
				c.Type = "enum"
				c.addToMore(match[1])
			}
		}
	}

	c.getMore(gormTab, adminTag)
}

func (c *OtherColumnInfo) getMore(gormTab string, adminTag string) {
	if strings.Index(gormTab, "not null") != -1 {
		c.More = append(c.More, "not null")
	}

	if strings.Index(adminTag, "disabled") != -1 {
		c.More  = append(c.More, "disabled")
	}
}


type BaseModel struct {
	Model interface{}
	Records interface{}
}

func NewBaseModel(model interface{}, records interface{}) BaseModel {
	return BaseModel{Model: model, Records: records}
}

func (m BaseModel) GetRecordAmount() int64 {
	var count int64
	adminDB.Model(m.Model).Count(&count)
	return count
}

func (m BaseModel) GetColumnInfo() []ColumnInfoDetail {

	var columnsInfo []ColumnInfoDetail

	re, _ := regexp.Compile(`comment:'(.+)'`)
	foreignKeyRe, _ := regexp.Compile(`foreignKey`)

	t := reflect.TypeOf(m.Model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}

	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		fmt.Println(t.Field(i).Type.Kind())
		name := t.Field(i).Name

		gormTab := t.Field(i).Tag.Get("gorm")
		jsonTag := t.Field(i).Tag.Get("json")
		adminTag := t.Field(i).Tag.Get("admin")
		other := OtherColumnInfo{}

		if foreignKeyRe.MatchString(gormTab) {
			if t.Field(i).Type.Kind() == reflect.Slice {

			} else {
				continue
			}
		}

		if jsonTag == "-" || adminTag == "-" {
			continue
		}

		other.GetInfo(adminTag, gormTab, t.Field(i).Type.Kind())

		if jsonTag == "" {
			jsonTag = name
		}

		var comment string
		r := re.FindStringSubmatch(gormTab)
		if len(r) > 1 {
			comment = r[1]
		} else {
			comment = jsonTag
		}

		//fmt.Println(adminTag, "===============")

		columnsInfo = append(columnsInfo, ColumnInfoDetail{jsonTag, comment, other})
	}

	return columnsInfo

}

func (m BaseModel) GetRecord(c *gin.Context) (gin.H, error) {
	fmt.Println("====admin get====")
	filter := NewFilter(m.Model)
	err := c.ShouldBindJSON(&filter)
	var count int64
	if err != nil {
		fmt.Println(err)
	} else {
		err := filter.SetFilter(&count)
		if err != nil {
			return gin.H{}, errors.New("find err")
		}
	}
	//fmt.Println(filter)

	result := filter.Finder.Order(`id desc`).Find(m.Records)

	if result.Error != nil {
		return gin.H{}, errors.New("find err")
	}

	return gin.H{
		"records":    m.Records,
		"columnInfo": m.GetColumnInfo(),
		"amount":     count,
	}, nil

}

func (m BaseModel) Updates(c *gin.Context) {
	err := c.ShouldBindJSON(m.Model)
	if err != nil {
		fmt.Println("ShouldBindJSON:", err)
		c.String(500, "ShouldBindJSON err")
		c.Abort()
		return
	}

	result := adminDB.Model(m.Model).Updates(m.Model)
	if result.Error != nil {
		c.String(500, "update err")
	} else {
		c.String(200, "update ok")
	}

}

func (m BaseModel) Create(c *gin.Context) {
	err := c.ShouldBindJSON(m.Model)
	if err != nil {
		fmt.Println(err)
		c.String(500, "ShouldBindJSON err")
		c.Abort()
		return
	}
	v := reflect.ValueOf(m.Model)
	f := v.Elem().FieldByName("ID")
	f.SetUint(0)

	result := adminDB.Create(m.Model)
	if result.Error != nil {
		c.String(500, "create err")
	} else {
		c.String(200, "create ok")
	}
}

func (m BaseModel) Delete(c *gin.Context) {
	bindStruct := struct {
		Ids []int `json:"ids"`
	}{}

	err := c.ShouldBindJSON(&bindStruct)
	if err != nil {
		fmt.Println("delete err:", err.Error())
		c.String(400, "bind ids err")
		return
	}

	result := adminDB.Delete(m.Records, bindStruct.Ids)
	if result.Error != nil {
		c.String(400, "delete err")
		return
	}

	c.String(200, "delete ok")
}

