package frame

import "testing"

func TestBrowserSetDefaults(t *testing.T) {
	browser := &Browser{}
	browser.SetDefaults()

	if browser.Executable == "" {
		t.Fatalf("browser executable not set")
	}
	if len(browser.Args) == 0 {
		t.Fatalf("browser args not set")
	}
	if browser.URL == "" {
		t.Fatalf("browser URL not set")
	}
}

func TestStartBrowser(t *testing.T) {
	browser := &Browser{}
	browser.SetDefaults()

	dummy := make(chan bool)
	browser.startBrowser(dummy)
}

func TestStartSlideshow(t *testing.T) {
	browser := &Browser{}
	browser.SetDefaults()

	dummy := make(chan bool)
	go browser.startBrowser(dummy)

	browser.startSlideshow(10)
}

func TestRunForever(t *testing.T) {
	browser := &Browser{}
	browser.SetDefaults()

	browser.RunForever()
}
