package emailcontroller

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type Mailaddress struct {
	To string `form:"to" json:"to"`
}

func OtpMail(c *gin.Context) {

	//OTP code
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	o := r1.Intn(999999)

	if o < 100000 {
		o += 100000
	}
	otp := strconv.Itoa(o)

	var address Mailaddress
	if err := c.ShouldBind(&address); err != nil {
		log.Fatal(err)
	}

	mailjetClient := mailjet.NewMailjetClient("9fdac63da0b2be01ee020adc6e42f4f8", "eee874c1379e2763460141cd800ce148")
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "teamsix756@gmail.com",
				Name:  "teamsix",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: address.To,
					Name:  address.To,
				},
			},
			Subject:  "OTP",
			TextPart: "Your six digit otp is :" + otp,
			HTMLPart: "<h3>Your otp is :</h3>" + otp,
			CustomID: "AppGettingStartedTest",
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusOK, gin.H{
			"status": "failed to send the otp",
			"error":  err,
		})
		return
	}
	fmt.Println(res)
	c.JSON(http.StatusOK, gin.H{
		"status": "varification opt sent",
		"otp":    otp,
	})
}
