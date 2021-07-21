//Package metadata 存放项目的一些metadata信息
// 这些参数可以在build 的时候通过 -ldflags的方式来设定
package metadata

import (
	"fmt"
)

var (
	// Name 项目名称
	Name = "Unknown"
	// Version 项目版本号
	Version = "Unknown"
	// 项目的平台
	Platform = "Unknown"
)

const versionString = `
%s (%s) on %s

---- Time is Money My Friend -----

`

func ShowVersion() {
	fmt.Printf(versionString, Name, Version, Platform)
}
