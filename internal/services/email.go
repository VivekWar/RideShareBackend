package services

import (
    "fmt"
    "rideshare-backend/internal/config"
    
    "gopkg.in/gomail.v2"
)

type EmailService struct {
    config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
    return &EmailService{
        config: cfg,
    }
}

func (s *EmailService) SendWelcomeEmail(email, name string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", s.config.EmailUser)
    m.SetHeader("To", email)
    m.SetHeader("Subject", "Welcome to RideShare!")
    
    body := fmt.Sprintf(`
        <html>
        <body>
            <h2>Welcome to RideShare, %s!</h2>
            <p>Thank you for joining our ride-sharing community.</p>
            <p>You can now:</p>
            <ul>
                <li>Create trips and find passengers</li>
                <li>Search for rides and join other travelers</li>
                <li>Save money on transportation costs</li>
                <li>Meet new people and make connections</li>
            </ul>
            <p>Get started by visiting your dashboard: <a href="%s/dashboard">%s/dashboard</a></p>
            <p>Happy travels!</p>
            <p>The RideShare Team</p>
        </body>
        </html>
    `, name, s.config.FrontendURL, s.config.FrontendURL)
    
    m.SetBody("text/html", body)
    
    d := gomail.NewDialer(s.config.EmailHost, s.config.EmailPort, s.config.EmailUser, s.config.EmailPassword)
    
    if err := d.DialAndSend(m); err != nil {
        return fmt.Errorf("failed to send welcome email: %w", err)
    }
    
    return nil
}

func (s *EmailService) SendTripMatchNotification(passengerEmail, passengerName, driverName, from, to string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", s.config.EmailUser)
    m.SetHeader("To", passengerEmail)
    m.SetHeader("Subject", "Trip Match Found!")
    
    body := fmt.Sprintf(`
        <html>
        <body>
            <h2>Great news, %s!</h2>
            <p>You've successfully joined a trip:</p>
            <div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 15px 0;">
                <p><strong>Route:</strong> %s → %s</p>
                <p><strong>Driver:</strong> %s</p>
            </div>
            <p>Please check your dashboard for more details and contact information.</p>
            <p><a href="%s/dashboard/trips" style="background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">View Trip Details</a></p>
            <p>Safe travels!</p>
            <p>The RideShare Team</p>
        </body>
        </html>
    `, passengerName, from, to, driverName, s.config.FrontendURL)
    
    m.SetBody("text/html", body)
    
    d := gomail.NewDialer(s.config.EmailHost, s.config.EmailPort, s.config.EmailUser, s.config.EmailPassword)
    
    return d.DialAndSend(m)
}

func (s *EmailService) SendTripCancellationEmail(email, name, from, to string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", s.config.EmailUser)
    m.SetHeader("To", email)
    m.SetHeader("Subject", "Trip Cancelled")
    
    body := fmt.Sprintf(`
        <html>
        <body>
            <h2>Trip Cancellation Notice</h2>
            <p>Dear %s,</p>
            <p>We regret to inform you that the following trip has been cancelled:</p>
            <div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 15px 0;">
                <p><strong>Route:</strong> %s → %s</p>
            </div>
            <p>We apologize for any inconvenience caused. Please search for alternative trips on our platform.</p>
            <p><a href="%s/search" style="background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Find Alternative Trips</a></p>
            <p>Thank you for your understanding.</p>
            <p>The RideShare Team</p>
        </body>
        </html>
    `, name, from, to, s.config.FrontendURL)
    
    m.SetBody("text/html", body)
    
    d := gomail.NewDialer(s.config.EmailHost, s.config.EmailPort, s.config.EmailUser, s.config.EmailPassword)
    
    return d.DialAndSend(m)
}
