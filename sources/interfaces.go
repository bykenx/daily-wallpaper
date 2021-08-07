package sources

type Source interface {
	// GetToday 获取今日图片
	GetToday() (TodayResponse, error)
	// GetArchive 获取历史图片
	GetArchive(param ArchiveParam) (ArchiveResponse, error)
}
