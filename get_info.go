package gologging

import (
	"fmt"
	"github.com/Laky-64/gologging/types"
	"path"
	"runtime"
	"strconv"
	"strings"
)

var (
	goRoutineFunction = fmt.Errorf("goroutine found")
	runtimeFunction   = fmt.Errorf("runtime function")
	testFunction      = fmt.Errorf("testing function")
)

func getInfo(skips int) (*types.CallerInfo, error) {
	var callerInfo types.CallerInfo
	pc, file, line, _ := runtime.Caller(skips + 1)
	callerInfo.Line = line
	callerInfo.FileName = path.Base(file)
	callerInfo.FuncName = runtime.FuncForPC(pc).Name()
	if strings.HasPrefix(callerInfo.FuncName, "runtime.goexit") {
		return nil, goRoutineFunction
	} else if strings.HasPrefix(callerInfo.FuncName, "runtime.") {
		return nil, runtimeFunction
	} else if strings.HasPrefix(callerInfo.FuncName, "testing.") {
		return nil, testFunction
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
	callerInfo.FuncName = funcInfo[5]
	if lambdaMatches := lambdaNameRgx.FindAllStringSubmatch(callerInfo.FuncName, -1); len(lambdaMatches) > 0 {
		numFunc, _ := strconv.Atoi(lambdaMatches[0][3])
		callerInfo.FuncName = fmt.Sprintf("lambda$%s$%d", lambdaMatches[0][2], numFunc-1)
	}
	callerInfo.FuncName = strings.ReplaceAll(callerInfo.FuncName, ".", "")
	return &callerInfo, nil
}
