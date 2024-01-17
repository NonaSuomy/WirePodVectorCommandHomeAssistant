/* 
Plugin: WirePod Vector Command HomeAssistant
Description: Sends the command after the word assist to HomeAssistant Conversation API "Hey Vector" wait for the ding... "Assist turn off family room lights"
Author: NonaSuomy
Date: 2024-01-17
*/

package main

import (
        "bytes"
        "encoding/json"
        "io/ioutil"
        "net/http"
        "strings"
)

var Utterances = []string{"assist"} // Trigger word
var Name = "Home Assistant Control"

type RequestBody struct {
        //AgentID string `json:"agent_id"`
        Text string `json:"text"`
}

type ResponseBody struct {
        Response struct {
                Speech struct {
                        Plain struct {
                                Speech string `json:"speech"`
                        } `json:"plain"`
                } `json:"speech"`
        } `json:"response"`
}
func stripOutTriggerWords(s string) string {
        result := strings.Replace(s, "assist", "", 1)
        return strings.TrimSpace(result) // remove leading and trailing spaces
}

func Action(transcribedText string, botSerial string, guid string, target string) (string, string) {
        if !strings.HasPrefix(transcribedText, "assist") {
                return "", "" // If the text does not start with "assist", do nothing
        }

        command := stripOutTriggerWords(transcribedText) // Get the command after "assist"

        url := "http://HAIPHERE:8123/api/conversation/process" // Replace with your Home Assistant IP
        token := "LongLiveTokenHere" // Replace with your token
        //agentID := "AgentIDHere" // Replace with your agent_id (Can get this with the dev assist console in yaml view or try the name)
        
        requestBody := &RequestBody{
                //AgentID: agentID,
                Text: command,
        }
        jsonValue, _ := json.Marshal(requestBody)

        req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
        req.Header.Set("Authorization", "Bearer "+token)
        req.Header.Set("Content-Type", "application/json")

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
                panic(err)
        }
        defer resp.Body.Close()

        body, _ := ioutil.ReadAll(resp.Body)

        var responseBody ResponseBody
        json.Unmarshal(body, &responseBody)

        VECTOR_PHRASE := responseBody.Response.Speech.Plain.Speech

        return "intent_imperative_praise", VECTOR_PHRASE
}
