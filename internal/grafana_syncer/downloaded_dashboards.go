package grafana_syncer

func NewDownloadedDashboards() DownloadedDashboards {
	var files []string
	return DownloadedDashboards{
		files: &files,
	}
}

type DownloadedDashboards struct {
	files *[]string
}

func (u DownloadedDashboards) markAsDownloaded(dashboard Dashboard) {
	*u.files = append(*u.files, dashboard.filename)
}

func (u DownloadedDashboards) hasBeenDownloaded(filename string) bool {
	for _, updatedFile := range *u.files {
		if updatedFile == filename {
			return true
		}
	}
	return false
}
