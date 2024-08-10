package gologging

import (
	"fmt"
	"github.com/Laky-64/gologging/types"
	"path"
	"runtime"
	"strconv"
	"strings"
)

func getInfo(skips int) (*types.CallerInfo, error) {
	var callerInfo types.CallerInfo
	pc, file, line, _ := runtime.Caller(skips + 1)
	callerInfo.Line = line
	callerInfo.FileName = path.Base(file)
	callerInfo.FuncName = runtime.FuncForPC(pc).Name()
	if strings.HasPrefix(callerInfo.FuncName, "runtime.") {
		return nil, fmt.Errorf("runtime function")
	} else if strings.HasPrefix(callerInfo.FuncName, "testing.") {
		return nil, fmt.Errorf("testing function")
	}
	funcInfo := getFunctionInfoRgx.FindStringSubmatch(callerInfo.FuncName)
	callerInfo.PackageName = strings.ToLower(
		strings.ReplaceAll(
			strings.ReplaceAll(
				funcInfo[1],
				"/",
				".",
			),
			"-",
			"_",
		),
	)
	callerInfo.FilePath = path.Join(path.Join(strings.Split(funcInfo[1], "/")[1:]...), path.Base(file))
	callerInfo.FuncName = funcInfo[3]
	if lambdaMatches := lambdaNameRgx.FindAllStringSubmatch(callerInfo.FuncName, -1); len(lambdaMatches) > 0 {
		lambdaDetails, err := getInfo(skips + 2)
		if err != nil {
			return nil, err
		}
		numFunc, _ := strconv.Atoi(lambdaMatches[0][2])
		callerInfo.FuncName = fmt.Sprintf("lambda$%s$%d", lambdaDetails.FuncName, numFunc-1)
	}
	callerInfo.FuncName = strings.ReplaceAll(callerInfo.FuncName, ".", "")
	return &callerInfo, nil
}
