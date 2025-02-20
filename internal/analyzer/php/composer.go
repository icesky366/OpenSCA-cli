/*
 * @Descripation: composer文件
 * @Date: 2021-11-26 14:50:06
 */

package php

import (
	"encoding/json"
	"opensca/internal/logs"
	"opensca/internal/srt"
)

// composer.lock
type ComposerLock struct {
	Pkgs []struct {
		Name    string            `json:"name"`
		Version string            `json:"version"`
		Require map[string]string `json:"require"`
	} `json:"packages"`
}

/**
 * @description: 解析composer.lock文件
 * @param {*srt.DepTree} depRoot 依赖树节点
 * @param {*srt.FileData} file 文件数据
 * @return {[]*srt.DepTree} 组件依赖列表
 */
func parseComposerLock(depRoot *srt.DepTree, file *srt.FileData) (deps []*srt.DepTree) {
	deps = []*srt.DepTree{}
	lock := ComposerLock{}
	if err := json.Unmarshal(file.Data, &lock); err != nil {
		logs.Error(err)
		return
	}
	// 记录组件信息
	depMap := map[string]*srt.DepTree{}
	for _, cps := range lock.Pkgs {
		dep := srt.NewDepTree(nil)
		dep.Name = cps.Name
		dep.Version = srt.NewVersion(cps.Version)
		depMap[cps.Name] = dep
	}
	// 构建依赖树
	for _, cps := range lock.Pkgs {
		for n := range cps.Require {
			if sub, ok := depMap[n]; ok && sub.Parent == nil {
				dep := depMap[cps.Name]
				sub.Parent = dep
				dep.Children = append(dep.Children, sub)
			}
		}
	}
	// 将顶层节点迁移到根节点下
	for _, cps := range lock.Pkgs {
		dep := depMap[cps.Name]
		if dep.Parent == nil {
			dep.Parent = depRoot
			depRoot.Children = append(depRoot.Children, dep)
		}
		deps = append(deps, dep)
	}
	return
}
