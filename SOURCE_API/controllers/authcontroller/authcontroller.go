package authcontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/team_six/SOURCE_API/models"

	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
)

var Entries []models.UserObject

func Auth(c *gin.Context) {
	var loginDetails models.Login
	if err := c.BindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Auth function is running")
	baseDN := "DC=mkclindia, DC=Local"
	ldapServer := "10.1.70.104:389"
	filterDN := "(&(objectClass=*)(sAMAccountName={username}))"
	ldapUsername := "mkclindia\\ProNExTTest"
	ldapPassword := "Mkcl#5050#"
	ldapConnection, ldapConnectionError := ldap.Dial("tcp", ldapServer)
	if ldapConnectionError != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": ldapConnectionError.Error(),
		})
		fmt.Println("ldapmdl connection Error: ", ldapConnectionError)
		return
	}
	fmt.Println(ldapConnection)
	ldapBindError := ldapConnection.Bind(ldapUsername, ldapPassword)
	if ldapBindError != nil {
		fmt.Println("ldapmdl ldapBindError: ", ldapBindError)
		return
	}
	result, searchError := ldapConnection.Search(ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter(loginDetails.Username, filterDN),
		[]string{"dn", "sAMAccountName", "mail", "sn", "givenName", "mobile"},
		nil,
	))
	if searchError != nil {
		fmt.Println("ldapmcl searcherror: ", searchError)
		c.JSON(http.StatusNotFound, gin.H{"status": searchError})
		return
	}
	if len(result.Entries) < 1 {
		fmt.Println("User does not exists")
		c.JSON(http.StatusNotFound, gin.H{"status": "User not found"})
		return
	}
	if len(result.Entries) > 1 {
		fmt.Println("Multiple entries of same UserID")
		c.JSON(http.StatusNotFound, gin.H{"status": "multiple entries of same user"})

		return
	}
	if userCredentialsBindError := ldapConnection.Bind(result.Entries[0].DN, loginDetails.Password); userCredentialsBindError != nil {
		fmt.Println("ldapmdl: Invalid Credentials")
		c.JSON(http.StatusNotFound, gin.H{"status": "Invalid Credentials"})
		return
	}
	var u models.UserObject
	u.FirstName = result.Entries[0].GetAttributeValue("givenName")
	u.LastName = result.Entries[0].GetAttributeValue("sn")
	u.MobileNumber = result.Entries[0].GetAttributeValue("mobile")
	u.Email = result.Entries[0].GetAttributeValue("mail")
	u.Username = result.Entries[0].GetAttributeValue("sAMAccountName")
	fmt.Println(u)
	c.JSON(200, gin.H{
		"status": "Successfully logged in",
	})
}
func filter(needle string, filterDN string) string {
	res := strings.Replace(filterDN, "{username}", needle, -1)
	return res

}
