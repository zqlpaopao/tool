package src

import (
	"github.com/zqlpaopao/tool/string-byte/src"
	"os/exec"
	"strings"
)

//getBranchInfo
//85e24f6810252fb50ca40110dda4ded7075cc0fb 85e24f6 6 天前 2022-01-04 20:34:45 +0800 1641299685 terinal zhangqiuli24
func getCmdInfo(command string,args... string) (info []string, err error){
	var out []byte
	cmd := exec.Command(command,args...)
	if out, err = cmd.Output();nil != err{
		return
	}

	return strings.Split(strings.Trim(src.Bytes2String(out),"\r\n")," "),nil
}
