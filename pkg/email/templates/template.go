package templates

import (
	"bytes"
	"go.uber.org/zap"
	"text/template"
	"time"
)

// Template represents an email template
type Template struct {
	Name       string      // Name of the template
	ContentRaw string      // Raw HTML content with potential template variables
	Data       interface{} // Data containing the template variables values
}

// Render renders a template with the given data
func (t Template) Render() (string, error) {

	// Render the parent template
	tmpl, err := template.New(t.Name).Parse(t.ContentRaw)
	if err != nil {
		return "", err
	}

	// Execute the template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, t.Data)
	if err != nil {
		return "", err
	}

	// Return the rendered template
	return buf.String(), nil
}

// Build renders the template and wraps it in the layout template
func (t Template) Build() (string, error) {

	// Render content
	content, err := t.Render()
	if err != nil {
		zap.L().Error("Render otp template", zap.Error(err))
		return "", err
	}

	// Render content within layout template
	emailTemplate := newLayoutTemplate(content)
	emailHtmlContent, err := emailTemplate.Render()
	if err != nil {
		zap.L().Error("Render email template", zap.Error(err))
		return "", err
	}

	return emailHtmlContent, nil
}

func newLayoutTemplate(content string) Template {
	return Template{
		Name:       "layout",
		ContentRaw: layoutHtml,
		Data: layoutData{
			Css:     layoutCss,
			Content: content,
			Year:    time.Now().Year(),
		},
	}
}

type layoutData struct {
	Css     string // CSS layout
	Content string // Content to be rendered within the layout
	Year    int    // Current year
}

const layoutHtml = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
	{{.Css}}
  </head>
  <body>
    <div class="body-wrapper">
      <main>
        <div class="content">
          {{.Content}}
        </div>
      </main>
      <p class="help">
        Need help? Contact us at
        <a href="mailto:contact@fihub.com">contact@fihub.com</a>
      </p>
      <footer>
        <h1>
          Fihub Company
        </h1>
        <p class="copyrights">
          Copyright © {{.Year}}. All rights reserved.
        </p>
      </footer>
    </div>
  </body>
</html>
`

const layoutCss = `
<link
	href="https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;500;600&display=swap"
	rel="stylesheet"
/>
<style>
      body {
        margin: 0;
        font-family: "Poppins", sans-serif;
        background: #ffffff;
        font-size: 14px;
      }
      .body-wrapper {
        max-width: 680px;
        margin: 0 auto;
        padding: 45px 30px 60px;
        font-size: 14px;
        color: #434343;
      }
      main {
        margin: 0;
        margin-top: 70px;
        padding: 50px 30px 15px;
        background: #ffffff;
        border-radius: 30px;
        text-align: center;
      }
      .content {
        width: 100%;
        max-width: 489px;
        margin: 0 auto;
      }
      .content h1 {
        margin: 0;
        font-size: 24px;
        font-weight: 500;
        color: #1f1f1f;
      }
      .content p {
        margin: 0;
        margin-top: 2rem;
        font-weight: 500;
        letter-spacing: 0.56px;
      }
      .content span.out {
        font-weight: 600;
        color: #1f1f1f;
      }
      .content p.secondary {
        margin-top: 1rem;
		color: #8c8c8c;
      }
      .content .otp {
        margin: 0;
        margin-top: 2rem;
        font-size: 40px;
        font-weight: 600;
        color: #8183f4;
        text-align: center;
      }
      .help {
        max-width: 400px;
        margin: 0 auto;
        margin-top: 50px;
        text-align: center;
        font-weight: 500;
        color: #8c8c8c;
      }
      footer {
        width: 100%;
        max-width: 490px;
        margin: 20px auto 0;
        text-align: center;
        border-top: 1px solid #e6ebf1;
      }
      footer h1 {
        margin: 0;
        margin-top: 1.5rem;
        font-size: 16px;
        font-weight: 600;
        color: #434343;
      }
      footer .copyrights {
        margin: 0;
        margin-top: 0.75rem;
        color: #434343;
      }
</style>
`
