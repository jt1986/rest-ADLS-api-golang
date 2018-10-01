package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	// working code below - please do not delete
	// Azure login credentials
	// cli cmd 'az login --serivce-principal -u ARM_CLIENT_ID -p ARM_CLIENT_SECRET --tenant ARM_TENANT_ID'
	ARM_RESOURCE := `https://management.core.windows.net`
	ARM_TENANT_ID := ""
	ARM_ClIENT_ID := ``
	ARM_CLIENT_SECRET := ``

	// Datalake Storage account name this needs to be passed in as a variable from Event Hub
	// Lists available Data Lake Store accounts.
	// cli cmd 'az dls account list --resource-group(-g) '
	// Get the details of a Data Lake Store account.
	// cli cmd 'az dls account show --account ARM_STORAGE_NAME'
	STORAGE_NAME := ""

	// ADLS file path this needs to be passed in as a variable from Event Hub
	// cli cmd 'az dls fs list --account ARM_STORAGE_NAME --path ARM_FILE_PATH'
	// Display the access control list (ACL).
	// cli cmd 'az dls fs access show --account ARM_STORAGE_NAME --path ARM_FILE_PATH'
	DATALAKE_FILE_PATH := ""

	//authenicate the connection
	var data map[string]interface{}
	pbody := strings.NewReader(`grant_type=client_credentials&resource=` + ARM_RESOURCE + `/&client_id=` + ARM_ClIENT_ID + `&client_secret=` + ARM_CLIENT_SECRET)
	psreq, pserr := http.NewRequest("POST", "https://login.microsoftonline.com/"+ARM_TENANT_ID+"/oauth2/token", pbody)
	if pserr != nil {
		// handle err
		panic(pserr)
	}
	psreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	presp, perr := http.DefaultClient.Do(psreq)
	if perr != nil {
		// handle err
		panic(perr)
	}
	defer presp.Body.Close()
	postBody, _ := ioutil.ReadAll(presp.Body)
	fmt.Println(string([]byte(postBody)))

	jsonerr := json.Unmarshal([]byte(postBody), &data)
	if jsonerr != nil {
		panic(jsonerr)
	}
	//obtain token from above authentication
	tok := data["access_token"].(string)
	berar := "Bearer " + tok

	//access the file location using HTTP
	fileRead := "https://" + STORAGE_NAME + ".azuredatalakestore.net/webhdfs/v1/" + DATALAKE_FILE_PATH + "?op=OPEN&read=true"
	getreq, geterr := http.NewRequest("GET", fileRead, nil)
	if geterr != nil {
		// handle err
	}
	fmt.Println(getreq)
	getreq.Header.Add("Authorization", berar)
	getresp, gerr := http.DefaultClient.Do(getreq)
	if gerr != nil {
		// handle err
	}
	defer getresp.Body.Close()

	body, _ := ioutil.ReadAll(getresp.Body)
	fmt.Println(string([]byte(body)))
	serviceData := []byte(body)

	// writing the read data to a file
	err := ioutil.WriteFile("dat1.csv", serviceData, 0644)
	if err != nil {
		panic(err)
	}

}
