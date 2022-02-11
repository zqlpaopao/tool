package src

import (
	"errors"
	"github.com/zqlpaopao/tool/format/src"
	"strings"
)

var version string

type versionNumManager struct {
	option      *option
	err         error
	version     string
	versionInfo []string
}

//func init(){
//	initVersion()
//}
//
//func initVersion(){
//	if err := NewVersionNumManager(
//		WithNotAuth(false),
//		WithBranch(true),
//		WithPrint(true),
//		WithTag("<========== Version Info ==========> "),
//	).Do().Error();err != nil{
//		src.PrintRed(err.Error())
//	}
//}

//NewVersionNumManager get versionManager
func NewVersionNumManager(f ...Options) *versionNumManager {
	return &versionNumManager{option: NewOptions(f...), err: nil, version: "", versionInfo: []string{}}
}

// Do make version
func (v *versionNumManager) Do()*versionNumManager {
	if v.option.witBranch {
		v.getBranch()
	}
	v.tidyInfo()
	if v.option.print{
		v.print()
	}
	return v
}

// tidyInfo auth time msg hashTag
func (v *versionNumManager) tidyInfo() {
	if v.err != nil {
		return
	}
	var versionInfo []string
	if versionInfo, v.err = getCmdInfo(GitCommand, GitArgs...); nil != v.err {
		return
	}
	if len(versionInfo) < 10 {
		v.err = errors.New("git info is error,less then 10")
		return
	}
	if !v.option.notAuth {
		v.versionInfo = append(v.versionInfo, versionInfo[9])
	}
	v.versionInfo = append(v.versionInfo,
		[]string{
			versionInfo[4] + " " + versionInfo[5],
			versionInfo[8],
			versionInfo[1],
		}...,
	)
	version = strings.Join(v.versionInfo, "_")
}

//getBranch get git branch
func (v *versionNumManager) getBranch() {
	v.versionInfo, v.err = getCmdInfo(GitCommand, GitBranchArgs...)
	return
}

//print print version
func (v *versionNumManager) print() {
	header := "******************************** Version info *****************************************"
	end := "***************************************************************************************"
	info := "**           "+v.option.tag
	for i:= 0;i<len(header)-len(v.option.tag)-15;i++{
		info += " "
	}
	info += "**"
	src.PrintGreenNoTime(header+`
`+info+`
`+end)
}

//Error get error
func (v *versionNumManager)Error()error{
	return v.err
}

//GetVersion Get version number
func GetVersion()string{
	return version
}