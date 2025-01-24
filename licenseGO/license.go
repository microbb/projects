package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strconv"
)

type AuthResponse struct {
	Token string `json:"token"`
}

type LicenseResponse struct {
	Total      int `json:"total"`
	Active     int `json:"active"`
	Supervisor int `json:"supervisor"`
}

type License struct {
	Company struct {
		Privileged struct {
			Used    int `json:"used"`
			Max     int `json:"max"`
			Overcap int `json:"overcap"`
		} `json:"privileged"`
	} `json:"_company"`
}

func main() {
	values := map[string]string{
		"auth_login": "",
		"password":   "",
	}

	token := authorized(values)

	priv := getLicense(token)
	comp := getLicenseActive(token)

	sendDataToZabbix(comp, priv)

}

func authorized(creds map[string]string) string {
	data, err := json.Marshal(creds)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := http.Post("https://elma.ozna.ru/guard/login", "application/json",
		bytes.NewBuffer(data))

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // response body is []byte

	var result AuthResponse

	err = json.Unmarshal(body, &result)
	fmt.Println(resp.StatusCode)

	return result.Token
}

func getLicense(token string) License {

	resp, _ := http.NewRequest(
		http.MethodPost,
		"https://elma.ozna.ru/api/status/portal/licenses/quantity",
		bytes.NewBufferString("[\"_company\"]"),
	)

	resp.Header.Add("Token", token)
	resp.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(resp)
	if err != nil {
		fmt.Println("Не удалось выполнить запрос на получение лицензии ", err)
	}

	defer resp.Body.Close()

	result := License{}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, &result)

	fmt.Println(response.StatusCode)
	return result
}

func getLicenseActive(token string) LicenseResponse {
	resp, _ := http.Get("https://elma.ozna.ru/api/status/competitive")
	resp.Header.Add("Token", token)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	result := LicenseResponse{}

	err = json.Unmarshal(body, &result)
	fmt.Println(resp.StatusCode)

	return result

}

func sendDataToZabbix(competitive LicenseResponse, privilegend License) {

	values := map[string]int{
		"elma_license_str": privilegend.Company.Privileged.Used,
		"elma_license_act": competitive.Active,
	}

	for k, v := range values {
		cmd := exec.Command(
			"zabbix_sender",
			"-z",
			"172.16.0.157",
			"-s",
			"test",
			"-p",
			"10051",
			"-k",
			k,
			"-o",
			strconv.Itoa(v),
		)

		err := cmd.Run()
		if err != nil {
			return
		}
	}
}
