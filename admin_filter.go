package agin

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Filter struct {
	Finder   *gorm.DB      `binding:"-"`
	Select   []string      `json:"select"`
	Where    []WhereFilter `json:"where"` // code:eq:xxx
	Paginate Paginate      `json:"paginate"`
}

func NewFilter(model interface{}) Filter {
	f := Filter{}
	f.Finder = G.DB.Model(model)
	return f
}

type WhereFilter struct {
	Type   string `json:"type"`
	Column string `json:"column"`
	Info   string `json:"info"`
}

type Paginate struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func (p Paginate) Offset() int {
	if p.Size > 30 {
		p.Size = 30
	}
	return (p.Page - 1) * p.Size
}

// 需要先InitFilterModel初始化
func (f *Filter) SetFilter(count *int64) error {
	fmt.Println("set filter====", f)
	if len(f.Select) > 0 {
		f.Finder.Select(f.Select)
	}
	if len(f.Where) > 0 {
		for _, w := range f.Where {
			if w.Type == "=" {
				f.Finder.Where(map[string]interface{}{w.Column: w.Info})
			} else if w.Type == "!=" {
				f.Finder.Not(map[string]interface{}{w.Column: w.Info})
			}
		}
	}
	result := f.Finder.Count(count)
	if result.Error != nil {
		return errors.New("get count err")
	}

	if f.Paginate.Size != 0 {
		f.Finder.Limit(f.Paginate.Size).Offset(f.Paginate.Offset())
	}

	return nil

}
