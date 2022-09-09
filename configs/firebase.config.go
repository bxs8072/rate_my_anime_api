package configs

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func FirebaseApp() (*firebase.App, error) {
	wd, _ := os.Getwd()
	opt := option.WithCredentialsFile(wd + "/configs/google_service.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return app, nil
}
