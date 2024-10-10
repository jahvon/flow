package io

import "github.com/jahvon/tuikit/styles"

func Theme(name string) styles.Theme {
	theme := styles.EverforestTheme()
	themeFunc, ok := styles.AllThemes()[name]
	if ok {
		theme = themeFunc()
	}
	return theme
}
