package email

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"strings"

	"github.com/romacardozx/stori-card-challenge/internal/config"
	"github.com/romacardozx/stori-card-challenge/internal/transaction"
)

func SendSummaryEmail(cfg config.SMTPConfig, summary transaction.Summary) error {
	// Crear el cliente SMTP
	smtpClient, err := smtp.Dial(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer smtpClient.Close()

	// Autenticar con el servidor SMTP
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	if err := smtpClient.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate with SMTP server: %v", err)
	}

	// Leer el archivo HTML de la plantilla
	htmlTemplate, err := ioutil.ReadFile("internal/email/assets/template.html")
	if err != nil {
		return fmt.Errorf("failed to read HTML template file: %v", err)
	}

	// Leer el archivo CSS
	css, err := ioutil.ReadFile("internal/email/assets/styles.css")
	if err != nil {
		return fmt.Errorf("failed to read CSS file: %v", err)
	}

	// Insertar el CSS en el HTML
	htmlContent := strings.Replace(string(htmlTemplate), "<link rel=\"stylesheet\" href=\"styles.css\">", "<style>"+string(css)+"</style>", 1)

	// Parsear el HTML como una plantilla
	tmpl, err := template.New("email").Parse(htmlContent)
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %v", err)
	}

	// Renderizar la plantilla con los datos del resumen
	var body bytes.Buffer
	err = tmpl.Execute(&body, summary)
	if err != nil {
		return fmt.Errorf("failed to execute HTML template: %v", err)
	}

	// Crear el mensaje de correo electrónico multiparte
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Agregar el cuerpo del correo electrónico como parte HTML
	htmlPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": {"text/html; charset=UTF-8"},
	})
	if err != nil {
		return fmt.Errorf("failed to create HTML part: %v", err)
	}
	_, err = htmlPart.Write(body.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write HTML part: %v", err)
	}

	// Adjuntar el logotipo como una parte inline
	logoPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":        {"image/png"},
		"Content-Disposition": {"inline; filename=logo.png"},
		"Content-ID":          {"<logo>"},
	})
	if err != nil {
		return fmt.Errorf("failed to create logo part: %v", err)
	}
	file, err := os.Open("internal/email/assets/logo.png")
	if err != nil {
		return fmt.Errorf("failed to open logo file: %v", err)
	}
	defer file.Close()
	_, err = io.Copy(logoPart, file)
	if err != nil {
		return fmt.Errorf("failed to copy logo file: %v", err)
	}

	// Finalizar el escritor multiparte
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// Establecer los encabezados del correo electrónico
	headers := make(textproto.MIMEHeader)
	headers.Set("From", cfg.From)
	headers.Set("To", cfg.To)
	headers.Set("Subject", "Resumen de Transacciones")
	headers.Set("Content-Type", fmt.Sprintf("multipart/related; boundary=%s", writer.Boundary()))

	// Enviar el correo electrónico
	err = smtpClient.Mail(cfg.From)
	if err != nil {
		return fmt.Errorf("failed to set mail from: %v", err)
	}
	err = smtpClient.Rcpt(cfg.To)
	if err != nil {
		return fmt.Errorf("failed to set mail to: %v", err)
	}

	emailData, err := smtpClient.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %v", err)
	}
	defer emailData.Close()

	for key, values := range headers {
		for _, value := range values {
			_, err := fmt.Fprintf(emailData, "%s: %s\r\n", key, value)
			if err != nil {
				return fmt.Errorf("failed to write header: %v", err)
			}
		}
	}

	_, err = fmt.Fprint(emailData, "\r\n")
	if err != nil {
		return fmt.Errorf("failed to write empty line: %v", err)
	}

	_, err = emailData.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write email data: %v", err)
	}

	return nil
}
