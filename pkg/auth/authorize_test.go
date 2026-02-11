package auth

import (
	"encoding/json"
	"fmt"
	"testing"
)

const jsonStr = `
{
	"email": "E0099999@cdtp.com",
	"name": "Wuque Hua",
	"extended_fields": {
		"HRBP": "7447",
		"SERVICECOMP": "1000",
	},
	"username": "E0099999",
	"externalId": "99999",
	"phone_number": "13488888888"
}
`

func TestJsonUnmarshal(t *testing.T) {

	userInfoMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(jsonStr), &userInfoMap); err != nil {
		fmt.Println("Parse user info error")
	} else {
		fmt.Println(userInfoMap)
	}
}

func TestGetUserFromRequest(t *testing.T) {
	user := parseIdentity(jsonStr, &Identity{authenticated: true})
	fmt.Println(user)
}
