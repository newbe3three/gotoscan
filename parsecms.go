package gotoscan

import (
	"encoding/json"
	"sort"
)

type CmsFeature struct {
	Path    string `json:"path"`
	Option  string `json:"option"`
	Content string `json:"content"`
}

//按照cms具有特征多少来排序，特征值多的先判断，特征少的后判断。
type CmsSort struct {
	Name   string
	number int
}

//实现sort接口来完排序
type CmsSortList []CmsSort

func (cmslist CmsSortList) Swap(i, j int)      { cmslist[i], cmslist[j] = cmslist[j], cmslist[i] }
func (cmslist CmsSortList) Len() int           { return len(cmslist) }
func (cmslist CmsSortList) Less(i, j int) bool { return cmslist[i].number > cmslist[j].number }

func ParseCmsFeatureFromJson(data []byte) (map[string][]CmsFeature, CmsSortList, error) {
	var cmslist map[string][]CmsFeature
	if err := json.Unmarshal(data, &cmslist); err != nil {
		return nil, nil, err
	}

	cmsSortList := make(CmsSortList, len(cmslist))
	i := 0
	for k, v := range cmslist {
		cmsSortList[i] = CmsSort{k, len(v)}
		i++
	}
	sort.Sort(&cmsSortList)

	return cmslist, cmsSortList, nil
}
