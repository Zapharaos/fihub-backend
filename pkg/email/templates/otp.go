package templates

// OtpData contains the data for the OTP template
type OtpData struct {
	Duration       string
	RequestLabel   string
	ProcedureLabel string
	OTP            string
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
<h1>Hello</h1>
<p>
{{.RequestLabel}}. Use the following OTP to complete the procedure to {{.ProcedureLabel}}. OTP is valid for <span class="out">{{.Duration}}</span>.
</p>
<p class="secondary">
Do not share this code with others, including Fihub employees.
</p>
<div class="otp">
	{{.OTP}}
</div>
`
