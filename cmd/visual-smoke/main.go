package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"

	"realtek-connect/internal/web"
)

const defaultChromePath = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"

type viewport struct {
	name   string
	width  int64
	height int64
}

type pageCheck struct {
	Title           string `json:"title"`
	HeroHeading     string `json:"heroHeading"`
	BodyTextLength  int64  `json:"bodyTextLength"`
	ScrollWidth     int64  `json:"scrollWidth"`
	ViewportWidth   int64  `json:"viewportWidth"`
	HeroImageLoaded bool   `json:"heroImageLoaded"`
	HasOverflow     bool   `json:"hasOverflow"`
}

func main() {
	var (
		baseURL    = flag.String("base-url", "", "existing base URL to check; defaults to an in-process local server")
		chromePath = flag.String("chrome-path", "", "path to the Chrome executable")
		timeout    = flag.Duration("timeout", 45*time.Second, "overall timeout for the smoke check")
	)
	flag.Parse()

	root, err := repoRoot()
	if err != nil {
		fail(err)
	}

	targetURL := *baseURL
	cleanup := func() {}
	if targetURL == "" {
		targetURL, cleanup, err = startLocalServer(root)
		if err != nil {
			fail(err)
		}
	}
	defer cleanup()

	chromeExec, err := resolveChromePath(*chromePath)
	if err != nil {
		fail(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	allocatorCtx, cancelAllocator := chromedp.NewExecAllocator(
		ctx,
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath(chromeExec),
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("hide-scrollbars", false),
			chromedp.NoDefaultBrowserCheck,
			chromedp.NoFirstRun,
		)...,
	)
	defer cancelAllocator()

	browserCtx, cancelBrowser := chromedp.NewContext(allocatorCtx)
	defer cancelBrowser()

	viewports := []viewport{
		{name: "desktop", width: 1440, height: 1100},
		{name: "mobile", width: 390, height: 844},
	}

	fmt.Printf("Visual smoke checks against %s\n", targetURL)
	for _, vp := range viewports {
		result, err := checkHome(browserCtx, targetURL, vp)
		if err != nil {
			fail(fmt.Errorf("%s viewport failed: %w", vp.name, err))
		}
		fmt.Printf(
			"- %s ok: title=%q hero=%q viewport=%d scrollWidth=%d heroImageLoaded=%t\n",
			vp.name,
			result.Title,
			result.HeroHeading,
			result.ViewportWidth,
			result.ScrollWidth,
			result.HeroImageLoaded,
		)
	}
}

func checkHome(parent context.Context, baseURL string, vp viewport) (pageCheck, error) {
	tabCtx, cancel := chromedp.NewContext(parent)
	defer cancel()

	var result pageCheck
	script := `(() => {
		const root = document.documentElement;
		const body = document.body;
		const hero = document.querySelector('.hero h1');
		const heroImage = document.querySelector('.hero-visual img');
		return {
			title: document.title,
			heroHeading: hero ? hero.textContent.trim() : "",
			bodyTextLength: body ? body.innerText.trim().length : 0,
			scrollWidth: root ? root.scrollWidth : 0,
			viewportWidth: window.innerWidth,
			heroImageLoaded: !!heroImage && heroImage.complete && heroImage.naturalWidth > 0,
			hasOverflow: !!root && root.scrollWidth > window.innerWidth
		};
	})()`

	if err := chromedp.Run(
		tabCtx,
		chromedp.EmulateViewport(vp.width, vp.height),
		chromedp.Navigate(baseURL),
		chromedp.WaitVisible(`main#main-content`, chromedp.ByQuery),
		chromedp.WaitVisible(`.hero`, chromedp.ByQuery),
		chromedp.WaitVisible(`.hero-visual img`, chromedp.ByQuery),
		chromedp.Sleep(300*time.Millisecond),
		chromedp.Evaluate(script, &result),
	); err != nil {
		return pageCheck{}, err
	}

	if result.Title == "" {
		return result, errors.New("document title is empty")
	}
	if result.BodyTextLength == 0 {
		return result, errors.New("page body is empty")
	}
	if result.HeroHeading != "Realtek Connect+" {
		return result, fmt.Errorf("unexpected hero heading %q", result.HeroHeading)
	}
	if !result.HeroImageLoaded {
		return result, errors.New("hero image did not load")
	}
	if result.HasOverflow {
		return result, fmt.Errorf("horizontal overflow detected: scrollWidth=%d viewportWidth=%d", result.ScrollWidth, result.ViewportWidth)
	}

	return result, nil
}

func startLocalServer(root string) (string, func(), error) {
	server, err := web.NewServer(web.Config{
		TemplatesDir: filepath.Join(root, "templates"),
		StaticDir:    filepath.Join(root, "static"),
	})
	if err != nil {
		return "", nil, err
	}

	testServer := httptest.NewServer(server.Routes())
	cleanup := func() {
		testServer.Close()
	}
	return testServer.URL, cleanup, nil
}

func repoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("could not locate repo root from current working directory")
		}
		dir = parent
	}
}

func resolveChromePath(flagValue string) (string, error) {
	candidates := []string{
		flagValue,
		os.Getenv("CHROME_PATH"),
		defaultChromePath,
		"google-chrome",
		"chromium",
		"chromium-browser",
		"chrome",
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if filepath.IsAbs(candidate) {
			if _, err := os.Stat(candidate); err == nil {
				return candidate, nil
			}
			continue
		}
		if resolved, err := exec.LookPath(candidate); err == nil {
			return resolved, nil
		}
	}

	return "", errors.New("could not locate Chrome; pass -chrome-path or set CHROME_PATH")
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
