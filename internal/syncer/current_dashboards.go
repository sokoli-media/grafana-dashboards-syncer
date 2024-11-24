package syncer

func NewCurrentDashboards() *CurrentDashboards {
	fileNameToDashboard := make(map[string]Dashboard)
	return &CurrentDashboards{
		filenameToDashboard: fileNameToDashboard,
	}
}

type CurrentDashboards struct {
	filenameToDashboard map[string]Dashboard
}

func (c CurrentDashboards) saveDashboard(dashboard Dashboard) {
	c.filenameToDashboard[dashboard.filename] = dashboard
}

func (c CurrentDashboards) dashboardHasBeenUpdated(newDashboard Dashboard) bool {
	oldDashboard, exists := c.filenameToDashboard[newDashboard.filename]
	return !exists || !oldDashboard.Equals(newDashboard)
}
