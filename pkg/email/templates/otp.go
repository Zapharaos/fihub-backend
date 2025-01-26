package templates

// OtpData contains the data for the OTP template
type OtpData struct {
	OTP         string
	Greeting    string
	MainContent string
	DoNotShare  string
}

// NewOtpTemplate creates a new OTP template
func NewOtpTemplate(data OtpData) Template {
	// Prepare otp template
	return Template{
		Name:       "otp",
		ContentRaw: otpHtml,
		Data:       data,
	}
}

const otpHtml = `
<h1>
	{{.Greeting}}
</h1>
<p>
	{{.MainContent}}
</p>
<p class="secondary">
	{{.DoNotShare}}
</p>
<div class="otp">
	{{.OTP}}
</div>
`
