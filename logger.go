package progress

func Logger(to func(format string, v ...any), prefix string, r Reader) {
	sugar := Extend(r)
	to("%s  0%%", prefix)
	for sugar.InProgress() {
		to(`%s%3d%% (%.2f/s, %s remaining)`, prefix, int(100*sugar.Percentage()), sugar.PerSecond(), sugar.Remaining())
	}
	to("%s100%% (%.2f/s, %s remaining)", prefix, sugar.PerSecond(), sugar.Remaining())
}
