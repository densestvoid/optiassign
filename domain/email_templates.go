package domain

import (
	"fmt"
	"html/template"
	"strings"
)

// EmailTemplateService handles email template generation
type EmailTemplateService struct {
	baseURL string
}

// NewEmailTemplateService creates a new email template service
func NewEmailTemplateService(baseURL string) *EmailTemplateService {
	return &EmailTemplateService{baseURL: baseURL}
}

// InvitationEmailData represents data for invitation emails
type InvitationEmailData struct {
	GroupName        string
	ParticipantName  string
	InvitationURL    string
	OwnerName        string
	DistributionRule string
	ItemCount        int
	Items            []*Item
}

// AssignmentEmailData represents data for assignment emails
type AssignmentEmailData struct {
	GroupName        string
	ParticipantName  string
	AssignedItems    []*Item
	TotalItems       int
	ResultsURL       string
	OwnerName        string
}

// GenerateInvitationHTML generates HTML invitation email
func (s *EmailTemplateService) GenerateInvitationHTML(data *InvitationEmailData) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Invitation to {{.GroupName}}</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #f8f9fa; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
        .button { display: inline-block; background: #007bff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; margin: 20px 0; }
        .items { background: #f8f9fa; padding: 15px; border-radius: 4px; margin: 15px 0; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>You're invited to participate in: {{.GroupName}}</h2>
            <p>Hello! {{.OwnerName}} has invited you to participate in a group assignment.</p>
        </div>
        
        <h3>Assignment Details</h3>
        <p><strong>Distribution Rule:</strong> {{.DistributionRule}}</p>
        <p><strong>Items to Rank:</strong> {{.ItemCount}} items</p>
        
        <div class="items">
            <h4>Items in this assignment:</h4>
            <ul>
                {{range .Items}}
                <li><strong>{{.Name}}</strong>{{if .Description}} - {{.Description}}{{end}}</li>
                {{end}}
            </ul>
        </div>
        
        <p>To participate, please click the button below and rank all items in order of your preference:</p>
        
        <a href="{{.InvitationURL}}" class="button">Submit Your Priorities</a>
        
        <div class="footer">
            <p>This invitation is unique to you. Please do not share this link.</p>
            <p>If you have any questions, please contact {{.OwnerName}}.</p>
        </div>
    </div>
</body>
</html>
`
	
	t, _ := template.New("invitation").Parse(tmpl)
	var buf strings.Builder
	t.Execute(&buf, data)
	return buf.String()
}

// GenerateAssignmentHTML generates HTML assignment email
func (s *EmailTemplateService) GenerateAssignmentHTML(data *AssignmentEmailData) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Your Assignment Results - {{.GroupName}}</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #d4edda; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
        .assigned-items { background: #f8f9fa; padding: 15px; border-radius: 4px; margin: 15px 0; }
        .item { padding: 10px; margin: 5px 0; background: white; border-left: 4px solid #007bff; }
        .button { display: inline-block; background: #007bff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; margin: 20px 0; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>🎉 Assignment Complete!</h2>
            <p>Your assignment results for <strong>{{.GroupName}}</strong> are ready!</p>
        </div>
        
        <h3>Your Assigned Items ({{len .AssignedItems}} of {{.TotalItems}})</h3>
        
        <div class="assigned-items">
            {{range $index, $item := .AssignedItems}}
            <div class="item">
                <strong>{{add $index 1}}. {{$item.Name}}</strong>
                {{if $item.Description}}<br><em>{{$item.Description}}</em>{{end}}
            </div>
            {{end}}
        </div>
        
        <p>You can view the complete results and details by clicking the button below:</p>
        
        <a href="{{.ResultsURL}}" class="button">View Complete Results</a>
        
        <div class="footer">
            <p>Thank you for participating in this assignment!</p>
            <p>If you have any questions, please contact {{.OwnerName}}.</p>
        </div>
    </div>
</body>
</html>
`
	
	t, _ := template.New("assignment").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"len": func(items []*Item) int { return len(items) },
	}).Parse(tmpl)
	var buf strings.Builder
	t.Execute(&buf, data)
	return buf.String()
}