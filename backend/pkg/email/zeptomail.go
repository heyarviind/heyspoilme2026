package email

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ZeptoMailClient struct {
	apiKey      string
	fromEmail   string
	fromName    string
	frontendURL string
	httpClient  *http.Client
}

type ZeptoMailRequest struct {
	From    ZeptoMailAddress   `json:"from"`
	To      []ZeptoMailTo      `json:"to"`
	Subject string             `json:"subject"`
	HTMLBody string            `json:"htmlbody"`
}

type ZeptoMailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name,omitempty"`
}

type ZeptoMailTo struct {
	EmailAddress ZeptoMailAddress `json:"email_address"`
}

func NewZeptoMailClient(apiKey, fromEmail, fromName, frontendURL string) (*ZeptoMailClient, error) {
	if apiKey == "" {
		return nil, errors.New("ZeptoMail API key is required")
	}
	if fromEmail == "" {
		fromEmail = "noreply@heyspoilme.com"
	}
	if fromName == "" {
		fromName = "HeySpoilMe"
	}

	return &ZeptoMailClient{
		apiKey:      apiKey,
		fromEmail:   fromEmail,
		fromName:    fromName,
		frontendURL: frontendURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (c *ZeptoMailClient) SendVerificationEmail(toEmail, token string) error {
	verifyURL := fmt.Sprintf("%s/auth/verify-email?token=%s", c.frontendURL, token)

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email</title>
</head>
<body style="margin: 0; padding: 0; background-color: #0a0a0a; font-family: 'Montserrat', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;">
    <table role="presentation" style="width: 100%%; max-width: 600px; margin: 0 auto; padding: 40px 20px;">
        <tr>
            <td style="text-align: center; padding-bottom: 30px;">
                <h1 style="color: #ffffff; font-size: 28px; margin: 0; font-weight: 600;">HeySpoilMe</h1>
            </td>
        </tr>
        <tr>
            <td style="background: rgba(255, 255, 255, 0.05); border: 1px solid rgba(255, 255, 255, 0.1); padding: 40px;">
                <h2 style="color: #ffffff; font-size: 24px; margin: 0 0 20px 0; font-weight: 500;">Verify Your Email</h2>
                <p style="color: rgba(255, 255, 255, 0.7); font-size: 16px; line-height: 1.6; margin: 0 0 30px 0;">
                    Welcome to HeySpoilMe! Please verify your email address to unlock all features including messaging, browsing filters, and profile image uploads.
                </p>
                <table role="presentation" style="width: 100%%;">
                    <tr>
                        <td style="text-align: center;">
                            <a href="%s" style="display: inline-block; background: #ffffff; color: #000000; padding: 16px 40px; text-decoration: none; font-weight: 600; font-size: 16px;">
                                Verify Email
                            </a>
                        </td>
                    </tr>
                </table>
                <p style="color: rgba(255, 255, 255, 0.5); font-size: 14px; line-height: 1.6; margin: 30px 0 0 0;">
                    This link will expire in 24 hours. If you didn't create an account on HeySpoilMe, you can safely ignore this email.
                </p>
            </td>
        </tr>
        <tr>
            <td style="text-align: center; padding-top: 30px;">
                <p style="color: rgba(255, 255, 255, 0.4); font-size: 12px; margin: 0;">
                    © 2026 HeySpoilMe. All rights reserved.
                </p>
            </td>
        </tr>
    </table>
</body>
</html>
`, verifyURL)

	return c.sendEmail(toEmail, "Verify your HeySpoilMe email", htmlBody)
}

func (c *ZeptoMailClient) SendNewMessageNotification(toEmail, senderName, messagePreview string) error {
	messagesURL := fmt.Sprintf("%s/messages", c.frontendURL)
	
	// Truncate message preview
	if len(messagePreview) > 100 {
		messagePreview = messagePreview[:100] + "..."
	}
	
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Message</title>
</head>
<body style="margin: 0; padding: 0; background-color: #0a0a0a; font-family: 'Montserrat', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;">
    <table role="presentation" style="width: 100%%; max-width: 600px; margin: 0 auto; padding: 40px 20px;">
        <tr>
            <td style="text-align: center; padding-bottom: 30px;">
                <h1 style="color: #ffffff; font-size: 28px; margin: 0; font-weight: 600;">HeySpoilMe</h1>
            </td>
        </tr>
        <tr>
            <td style="background: rgba(255, 255, 255, 0.05); border: 1px solid rgba(255, 255, 255, 0.1); padding: 40px;">
                <h2 style="color: #ffffff; font-size: 24px; margin: 0 0 20px 0; font-weight: 500;">You have a new message</h2>
                <p style="color: rgba(255, 255, 255, 0.7); font-size: 16px; line-height: 1.6; margin: 0 0 15px 0;">
                    <strong style="color: #ffffff;">%s</strong> sent you a message:
                </p>
                <div style="background: rgba(255, 255, 255, 0.03); border-left: 3px solid rgba(255, 255, 255, 0.3); padding: 15px 20px; margin: 0 0 30px 0;">
                    <p style="color: rgba(255, 255, 255, 0.6); font-size: 14px; line-height: 1.6; margin: 0; font-style: italic;">
                        "%s"
                    </p>
                </div>
                <table role="presentation" style="width: 100%%;">
                    <tr>
                        <td style="text-align: center;">
                            <a href="%s" style="display: inline-block; background: #ffffff; color: #000000; padding: 16px 40px; text-decoration: none; font-weight: 600; font-size: 16px;">
                                View Message
                            </a>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td style="text-align: center; padding-top: 30px;">
                <p style="color: rgba(255, 255, 255, 0.4); font-size: 12px; margin: 0;">
                    © 2026 HeySpoilMe. All rights reserved.
                </p>
                <p style="color: rgba(255, 255, 255, 0.3); font-size: 11px; margin: 10px 0 0 0;">
                    You received this email because you have unread messages. 
                </p>
            </td>
        </tr>
    </table>
</body>
</html>
`, senderName, messagePreview, messagesURL)

	return c.sendEmail(toEmail, fmt.Sprintf("New message from %s on HeySpoilMe", senderName), htmlBody)
}

func (c *ZeptoMailClient) sendEmail(toEmail, subject, htmlBody string) error {
	payload := ZeptoMailRequest{
		From: ZeptoMailAddress{
			Address: c.fromEmail,
			Name:    c.fromName,
		},
		To: []ZeptoMailTo{
			{
				EmailAddress: ZeptoMailAddress{
					Address: toEmail,
				},
			},
		},
		Subject:  subject,
		HTMLBody: htmlBody,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.zeptomail.in/v1.1/email", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Zoho-enczapikey "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		log.Printf("[ZeptoMail] Error response: %v", errorResp)
		return fmt.Errorf("ZeptoMail API error: status %d", resp.StatusCode)
	}

	log.Printf("[ZeptoMail] Successfully sent email to %s", toEmail)
	return nil
}

