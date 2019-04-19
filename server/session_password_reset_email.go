package server

import (
    "fmt"
    //go get -u github.com/aws/aws-sdk-go
    "github.com/aws/aws-sdk-go/aws"
    awsSession "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/service/ses"
    "github.com/aws/aws-sdk-go/aws/awserr"
)

const (
    //The email body for recipients with non-HTML email clients.
    TextBody = "Use this link to generate your sinal password"
    
    // The character encoding for the email.
    CharSet = "UTF-8"
)

func generatePasswordResetEmailBody(resetToken string, config Config) (string) {
    FullResetUrl := config.GetEmail().PasswordRecoveryUrl + "?reset_token=" + resetToken
    HtmlBody :=  "<h1>Dear Signal Player</h1><p>We have received a request to reset your " + config.GetEmail().AppName + " password. " +
                "If you are the person who requested this, please use the following link " + 
                "to update the password within 24 hours</p>" +
                "<a href='" + FullResetUrl + "'>" + FullResetUrl + "</a>" +
                "<p>Cheers,<br/>Signal Game Team</p>"
    return HtmlBody
}

func sendPasswordResetEmail(Recipient string, resetToken string, config Config) {
    // Create a new session in the us-west-2 region.
    // Replace us-west-2 with the AWS Region you're using for Amazon SES.
    sess, err := awsSession.NewSession(
        &aws.Config{Region: aws.String(config.GetEmail().AWSRegion),
        Credentials: credentials.NewStaticCredentials(config.GetEmail().AWSAccessKeyID, config.GetEmail().AWSSecretAccessKey,""),
    })
    
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
                    Data:    aws.String(generatePasswordResetEmailBody(resetToken, config)),
                },
                Text: &ses.Content{
                    Charset: aws.String(CharSet),
                    Data:    aws.String(TextBody),
                },
            },
            Subject: &ses.Content{
                Charset: aws.String(CharSet),
                Data:    aws.String(config.GetEmail().PasswordEmailSubject),
            },
        },
        Source: aws.String(config.GetEmail().PasswordEmailSender),
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