package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

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
	priv := getLicense()
	comp := getLicenseActive()

	sendDataToZabbix(comp, priv)

}

func getLicense() License {

	resp, _ := http.NewRequest(
		http.MethodPost,
		"https://elma.ozna.ru/api/status/portal/licenses/quantity",
		bytes.NewBufferString("[\"_company\"]"),
	)

	resp.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(resp)
	if err != nil {
		log.Println("Не удалось выполнить запрос на получение лицензии ", err)
	}

	defer resp.Body.Close()

	result := License{}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(body, &result)

	log.Println(response.StatusCode)
	return result
}

func getLicenseActive() LicenseResponse {
	resp, _ := http.Get("https://elma.ozna.ru/api/status/competitive")

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	result := LicenseResponse{}

	err = json.Unmarshal(body, &result)
	log.Println(resp.StatusCode)

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
