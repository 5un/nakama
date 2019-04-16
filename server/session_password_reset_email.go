package server

import (
    "fmt"
    
    //go get -u github.com/aws/aws-sdk-go
    "github.com/aws/aws-sdk-go/aws"
    awsSession "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ses"
    "github.com/aws/aws-sdk-go/aws/awserr"
)

const (
    // Replace sender@example.com with your "From" address. 
    // This address must be verified with Amazon SES.
    Sender = "accounts@mailer.signalvideogame.com"

    // Specify a configuration set. To use a configuration
    // set, comment the next line and line 92.
    //ConfigurationSet = "ConfigSet"
    
    // The subject line for the email.
    Subject = "Reset your Signalvideogame.com Account Password"
        
    //The email body for recipients with non-HTML email clients.
    TextBody = "Use this link to generate your sinal password"
    
    ResetUrl = "https://www.signalvideogame.com/password_recovery.html?reset_token="

    // The character encoding for the email.
    CharSet = "UTF-8"
)

func generatePasswordResetEmailBody(resetToken string) (string) {
    FullResetUrl := ResetUrl + resetToken
    HtmlBody :=  "<h1>Dear Signal Player</h1><p>We have received a request to reset your signalvideogame.com password. " +
                "If you are the person who requested this, please use the following link " + 
                "to update the password within 24 hours</p>" +
                "<a href='" + FullResetUrl + "'>" + FullResetUrl + "</a>" +
                "<p>Cheers,<br/>Signal Game Team</p>"
    return HtmlBody
}

func sendPasswordResetEmail(Recipient string, resetToken string) {
    // Create a new session in the us-west-2 region.
    // Replace us-west-2 with the AWS Region you're using for Amazon SES.
    sess, err := awsSession.NewSessionWithOptions(
        awsSession.Options{
            Config: aws.Config{Region: aws.String("us-east-1")},
            Profile: "default",
        },
    )
    
    // Create an SES session.
    svc := ses.New(sess)
    
    // Assemble the email.
    input := &ses.SendEmailInput{
        Destination: &ses.Destination{
            CcAddresses: []*string{
            },
            ToAddresses: []*string{
                aws.String(Recipient),
            },
        },
        Message: &ses.Message{
            Body: &ses.Body{
                Html: &ses.Content{
                    Charset: aws.String(CharSet),
                    Data:    aws.String(generatePasswordResetEmailBody(resetToken)),
                },
                Text: &ses.Content{
                    Charset: aws.String(CharSet),
                    Data:    aws.String(TextBody),
                },
            },
            Subject: &ses.Content{
                Charset: aws.String(CharSet),
                Data:    aws.String(Subject),
            },
        },
        Source: aws.String(Sender),
            // Uncomment to use a configuration set
            //ConfigurationSetName: aws.String(ConfigurationSet),
    }

    // Attempt to send the email.
    result, err := svc.SendEmail(input)
    
    // Display error messages if they occur.
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case ses.ErrCodeMessageRejected:
                fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
            case ses.ErrCodeMailFromDomainNotVerifiedException:
                fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
            case ses.ErrCodeConfigurationSetDoesNotExistException:
                fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
            default:
                fmt.Println(aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            fmt.Println(err.Error())
        }
    
        return
    }
    
    fmt.Println("Email Sent to address: " + Recipient)
    fmt.Println(result)
}