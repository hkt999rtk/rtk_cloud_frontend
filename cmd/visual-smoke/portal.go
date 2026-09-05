package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

// Exercise real browser events against the selected test server. The only POST
// is an empty contact form, which must fail validation without storing a lead.
func checkPortalInteractions(ctx context.Context, baseURL string) error {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return err
	}
	if host := parsed.Hostname(); host != "localhost" && host != "127.0.0.1" && host != "::1" {
		return fmt.Errorf("Portal interaction checks require a loopback preview; refusing to submit a contact form to %s", host)
	}
	assert := func(expression, label string) chromedp.Action {
		return chromedp.ActionFunc(func(ctx context.Context) error {
			var ok bool
			if err := chromedp.Evaluate(expression, &ok).Do(ctx); err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("Portal interaction failed: %s", label)
			}
			return nil
		})
	}
	return chromedp.Run(ctx,
		chromedp.EmulateViewport(390, 844),
		chromedp.Navigate(baseURL+"/"),
		chromedp.WaitVisible(`[data-menu-toggle]`, chromedp.ByQuery),
		chromedp.Click(`[data-menu-toggle]`, chromedp.ByQuery),
		assert(`document.querySelector('[data-menu-toggle]').getAttribute('aria-expanded') === 'true' && document.querySelector('#site-navigation').contains(document.activeElement)`, "mobile menu opens and receives focus"),
		chromedp.KeyEvent(kb.Escape),
		assert(`document.activeElement.matches('[data-menu-toggle]') && document.querySelector('[data-menu-toggle]').getAttribute('aria-expanded') === 'false'`, "Escape closes menu and restores focus"),
		chromedp.ScrollIntoView(`[data-core-accordion="panel-ota"]`, chromedp.ByQuery),
		chromedp.Click(`[data-core-accordion="panel-ota"]`, chromedp.ByQuery),
		assert(`document.querySelector('[data-core-accordion="panel-ota"]').getAttribute('aria-expanded') === 'true' && document.querySelector('#panel-ota-content').getBoundingClientRect().height > 0`, "mobile capability accordion"),
		chromedp.EmulateViewport(1440, 1100),
		chromedp.Navigate(baseURL+"/"),
		chromedp.Focus(`#tab-provision`, chromedp.ByQuery),
		chromedp.KeyEvent(kb.ArrowRight),
		assert(`document.activeElement.id === 'tab-ota' && document.querySelector('#tab-ota').getAttribute('aria-selected') === 'true' && !document.querySelector('#panel-ota').hidden`, "keyboard capability tabs"),
		chromedp.Navigate(baseURL+"/zh-tw/manual/getting-started"),
		chromedp.WaitVisible(`[data-manual-toc] a`, chromedp.ByQuery),
		assert(`document.querySelector('[data-manual-toc] h2').textContent === '本頁內容' && document.querySelectorAll('[data-manual-toc] a').length === document.querySelectorAll('[data-manual-article] h2').length`, "localized manual contents"),
		chromedp.Click(`[data-manual-toc] a`, chromedp.ByQuery),
		assert(`document.activeElement.matches('[data-manual-article] h2') && !!location.hash`, "manual anchor and keyboard focus"),
		chromedp.Navigate(baseURL+"/contact"),
		chromedp.Click(`.form-card button[type="submit"]`, chromedp.ByQuery),
		chromedp.WaitVisible(`.form-error-summary`, chromedp.ByQuery),
		assert(`document.querySelector('.form-error-summary').getAttribute('role') === 'alert' && document.querySelectorAll('[aria-invalid="true"]').length > 0`, "accessible contact validation"),
	)
}
