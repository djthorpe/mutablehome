package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	mutablehome "github.com/djthorpe/mutablehome"
)

var (
	wg sync.WaitGroup
)

/////////////////////////////////////////////////////////////////////

func GetFolder(app gopi.App, args []string) (string, error) {
	if len(args) == 0 {
		if wd, err := os.Getwd(); err != nil {
			return "", err
		} else {
			return wd, nil
		}
	}
	if len(args) == 1 {
		return args[0], nil
	}

	// Return help
	return "", gopi.ErrHelp
}

func ScanFolders(app gopi.App, stop <-chan struct{}, folder string) {
	// Wait for scan to complete
	wg.Add(1)
	defer wg.Done()

	// Scan folders for playlists, every 10 seconds
	timer := time.NewTimer(500 * time.Millisecond)
FOR_LOOP:
	for {
		select {
		case <-timer.C:
			if err := ScanFolder(app, folder); err != nil {
				app.Log().Error(err)
			}
			timer.Reset(10 * time.Second)
		case <-stop:
			timer.Stop()
			break FOR_LOOP
		}
	}
}

func ScanFolder(app gopi.App, folder string) error {
	return filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		if path == folder {
			return nil
		}
		if relpath, err := filepath.Rel(folder, path); err != nil {
			return nil
		} else {
			fmt.Println(relpath, info.ModTime())
		}
		return nil
	})
}

func Main(app gopi.App, args []string) error {
	httpd := app.UnitInstance("httpd").(mutablehome.HttpServer)
	stop := make(chan struct{})

	// Serve HTTP
	if folder, err := GetFolder(app, args); err != nil {
		return err
	} else if url, err := httpd.ServeStatic(folder); err != nil {
		return err
	} else {
		fmt.Println("Serving on", url)
		go ScanFolders(app, stop, folder)
	}

	// Wait for CTRL+C
	fmt.Println("Press CTRL+C to end")
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Wait for end of scanning
	close(stop)
	wg.Wait()

	// Return success
	return nil
}
