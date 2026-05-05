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
	"strings"
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
	BodyText        string `json:"bodyText"`
	BodyTextLength  int64  `json:"bodyTextLength"`
	ScrollWidth     int64  `json:"scrollWidth"`
	ViewportWidth   int64  `json:"viewportWidth"`
	HeroImageLoaded bool   `json:"heroImageLoaded"`
	HasOverflow     bool   `json:"hasOverflow"`
}

type pageTarget struct {
	name             string
	path             string
	headingSelector  string
	imageSelector    string
	expectedHeading  string
	expectedTitle    string
	expectedBodyText string
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
	targets := []pageTarget{
		{
			name:             "home-en",
			path:             "/",
			headingSelector:  ".hero h1",
			imageSelector:    ".hero-visual img",
			expectedHeading:  "Realtek Connect+",
			expectedTitle:    "Realtek Connect+ | IoT Cloud Platform",
			expectedBodyText: "Realtek-based devices online",
		},
		{
			name:             "home-zh-tw",
			path:             "/zh-tw/",
			headingSelector:  ".hero h1",
			imageSelector:    ".hero-visual img",
			expectedHeading:  "Realtek Connect+",
			expectedTitle:    "Realtek Connect+ | 物聯網雲端平台",
			expectedBodyText: "讓 Realtek 裝置更快進入",
		},
		{
			name:             "home-zh-cn",
			path:             "/zh-cn/",
			headingSelector:  ".hero h1",
			imageSelector:    ".hero-visual img",
			expectedHeading:  "Realtek Connect+",
			expectedTitle:    "Realtek Connect+ | 物联网云端平台",
			expectedBodyText: "让 Realtek 装置更快进入",
		},
		{
			name:             "feature-zh-tw",
			path:             "/zh-tw/features/provision",
			headingSelector:  ".feature-detail h1",
			imageSelector:    ".feature-visual img",
			expectedHeading:  "Provision 配網",
			expectedTitle:    "Provision 配網 | Realtek Connect+",
			expectedBodyText: "合約支撐的基礎",
		},
		{
			name:             "feature-zh-cn",
			path:             "/zh-cn/features/provision",
			headingSelector:  ".feature-detail h1",
			imageSelector:    ".feature-visual img",
			expectedHeading:  "Provision 配网",
			expectedTitle:    "Provision 配网 | Realtek Connect+",
			expectedBodyText: "合约支撑的基础",
		},
		{
			name:             "manual-en",
			path:             "/manual/getting-started",
			headingSelector:  ".feature-detail h1",
			imageSelector:    ".manual-content img",
			expectedHeading:  "Getting Started",
			expectedTitle:    "Getting Started | Realtek Connect+",
			expectedBodyText: "Set up your first device",
		},
	}

	fmt.Printf("Visual smoke checks against %s\n", targetURL)
	for _, target := range targets {
		for _, vp := range viewports {
			result, err := checkPage(browserCtx, targetURL, target, vp)
			if err != nil {
				fail(fmt.Errorf("%s %s viewport failed: %w", target.name, vp.name, err))
			}
			fmt.Printf(
				"- %s %s ok: title=%q heading=%q viewport=%d scrollWidth=%d imageLoaded=%t\n",
				target.name,
				vp.name,
				result.Title,
				result.HeroHeading,
				result.ViewportWidth,
				result.ScrollWidth,
				result.HeroImageLoaded,
			)
		}
	}
}

func checkPage(parent context.Context, baseURL string, target pageTarget, vp viewport) (pageCheck, error) {
	tabCtx, cancel := chromedp.NewContext(parent)
	defer cancel()

	var result pageCheck
	script := fmt.Sprintf(`(() => {
		const root = document.documentElement;
		const body = document.body;
		const hero = document.querySelector(%q);
		const heroImage = document.querySelector(%q);
		return {
			title: document.title,
			heroHeading: hero ? hero.textContent.trim() : "",
			bodyText: body ? body.innerText.trim() : "",
			bodyTextLength: body ? body.innerText.trim().length : 0,
			scrollWidth: root ? root.scrollWidth : 0,
			viewportWidth: window.innerWidth,
			heroImageLoaded: !!heroImage && heroImage.complete && heroImage.naturalWidth > 0,
			hasOverflow: !!root && root.scrollWidth > window.innerWidth
		};
	})()`, target.headingSelector, target.imageSelector)

	if err := chromedp.Run(
		tabCtx,
		chromedp.EmulateViewport(vp.width, vp.height),
		chromedp.Navigate(baseURL+target.path),
		chromedp.WaitVisible(`main#main-content`, chromedp.ByQuery),
		chromedp.WaitVisible(target.headingSelector, chromedp.ByQuery),
		chromedp.WaitVisible(target.imageSelector, chromedp.ByQuery),
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
	if target.expectedTitle != "" && result.Title != target.expectedTitle {
		return result, fmt.Errorf("unexpected title %q", result.Title)
	}
	if result.HeroHeading != target.expectedHeading {
		return result, fmt.Errorf("unexpected heading %q", result.HeroHeading)
	}
	if target.expectedBodyText != "" && !containsText(result, target.expectedBodyText) {
		return result, fmt.Errorf("expected body text %q was not found", target.expectedBodyText)
	}
	if !result.HeroImageLoaded {
		return result, errors.New("target image did not load")
	}
	if result.HasOverflow {
		return result, fmt.Errorf("horizontal overflow detected: scrollWidth=%d viewportWidth=%d", result.ScrollWidth, result.ViewportWidth)
	}

	return result, nil
}

func containsText(result pageCheck, text string) bool {
	return result.BodyTextLength > int64(len(text)) && strings.Contains(result.BodyText, text)
}

func startLocalServer(root string) (string, func(), error) {
	server, err := web.NewServer(web.Config{
		TemplatesDir: filepath.Join(root, "templates"),
		StaticDir:    filepath.Join(root, "static"),
		ContentDir:   filepath.Join(root, "content", "docs"),
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
