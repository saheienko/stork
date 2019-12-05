package volumeimport

import "github.com/libopenstorage/stork/pkg/volumeimport/controllers"

type Controller struct {
}

func New() (*Controller, error) {
	return &Controller{}, nil
}

func (c Controller) Init() error {
	viController, err := controllers.NewVolumeImportController()
	if err != nil {
		return err
	}
	if err = viController.Init(); err != nil {
		return err
	}

	jobController, err := controllers.NewJobController()
	if err != nil {
		return err
	}
	if err = jobController.Init(); err != nil {
		return err
	}

	return nil
}
