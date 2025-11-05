package ismismcube_server

import (
	"ismismcube-backend/internal/model"
)

func GetPageViewCount() (int, error) {
	return model.GetIsmismcubePageViewCount()
}

func AddPageView(visitorIP, userAgent string) (*model.IsmismcubePageView, error) {
	view := &model.IsmismcubePageView{
		VisitorIP: visitorIP,
		UserAgent: userAgent,
	}
	return model.AddIsmismcubePageView(view)
}

func GetPageViews(limit, offset int) ([]model.IsmismcubePageView, error) {
	return model.GetIsmismcubePageViews(limit, offset)
}

func DeletePageView(id int) error {
	return model.DeleteIsmismcubePageView(id)
}
