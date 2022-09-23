package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reminders/app"
	"reminders/app/utils"
	"reminders/app/whatsapp"
	"reminders/app/workflows"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
	"go.temporal.io/sdk/client"
)

func (h *RequestHandler) ReminderListHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "RemindersList: %s\n", "")
}

func (h *RequestHandler) CreateReminderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
		return
	}

	var input utils.ReminderInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input.FromTime = time.Now()

	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	reminderInfo, err := workflows.StartWorkflow(c, &input)
	log.Printf("Creating reminder for Phone %s", input.Phone)
	if err != nil {
		log.Printf("failed to start workflow: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Created reminder for workflowId %s runId %s", reminderInfo.WorkflowId, reminderInfo.RunId)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	makeReminderResponse(w, reminderInfo)
}

func makeReminderResponse(w http.ResponseWriter, r utils.ReminderDetails) {
	json.NewEncoder(w).Encode(
		map[string]string{
			"ReferenceId":  r.ReferenceId,
			"ReminderName": r.ReminderName,
			"ReminderText": r.ReminderText,
			"ReminderTime": utils.GetReminderTime(r.FromTime, r.NMinutes).Format(app.TIME_FORMAT),
		})
}

func (h *RequestHandler) GetReminderHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"alive": true}`)
}

func (h *RequestHandler) UpdateReminderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	referenceId := vars["referenceId"]
	workflowId, runId, err := utils.GetInternalIdsFromReferenceId(referenceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input utils.ReminderInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	reminderInfo, err := workflows.UpdateWorkflow(c, workflowId, runId, &input)
	if err != nil {
		log.Printf("Failed to update workflow: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Updated reminder for workflowId %s runId %s", workflowId, runId)
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	makeReminderResponse(w, reminderInfo)
}

func (h *RequestHandler) DeleteReminderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	referenceId := vars["referenceId"]

	workflowId, runId, err := utils.GetInternalIdsFromReferenceId(referenceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := h.c.GetClient()
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer c.Close()

	err = workflows.DeleteWorkflow(c, whatsapp.GetWhatsappClient(), workflowId, runId)
	if err != nil {
		log.Printf("Failed to delete workflow %s (runID %s): %v", workflowId, runId, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Deleted reminder for workflowId %s runId %s", workflowId, runId)
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		map[string]string{
			"referenceId": referenceId,
			"status":      "cancelled",
		})
}

func (h *RequestHandler) WhatsappResponseHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("WhatsApp message received.")
	body, err := ioutil.ReadAll(r.Body)

	if r.Method == "GET" {
		handleVerification(w, r)
		return
	}

	if err != nil {
		http.Error(w, "Unable to read request body.", http.StatusBadRequest)
		return
	}

	results := gjson.GetManyBytes(
		body,
		"entry.0.changes.0.value.messages.0.from",
		"entry.0.changes.0.value.messages.0.timestamp",
		"entry.0.changes.0.value.messages.0.text.body",
	)

	fromPhone := results[0].Str
	timestampStr := results[1].Str
	message := results[2].Str

	if fromPhone == "" {
		http.Error(w, "From phone number not found in request.", http.StatusBadRequest)
		return
	}

	timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid WhatsApp request.", http.StatusBadRequest)
		return
	}
	fromTime := time.Unix(timestampInt, 0)

	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	reminderInfo, err := doMessageAction(c, fromPhone, message, fromTime)

	if err != nil {
		log.Print("Sending Whatsapp Error message")
		whatsapp.GetWhatsappClient().SendMessage(fromPhone, "Unable to create reminder; unrecognized request format.")
		http.Error(w, "Unrecognized reminder request format.", http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		makeReminderResponse(w, reminderInfo)
	}
}

func handleVerification(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query()
	verifyToken, tokenFound := queryString["hub.verify_token"]
	challenge, challengeFound := queryString["hub.challenge"]

	if tokenFound == false || challengeFound == false || verifyToken[0] != app.FB_VERIFY_TOKEN {
		http.Error(w, "Unrecognized request.", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fbResp := fmt.Sprintf("%s", challenge[0])
	io.WriteString(w, fbResp)
	return
}

func doMessageAction(c client.Client, phone string, message string, fromTime time.Time) (utils.ReminderDetails, error) {
	if name, text, nMinutes, err := app.ParseCreateReminderMessage(message); err == nil {
		return createReminderFromMessage(c, phone, name, text, nMinutes, fromTime)
	}
	if referenceId, nMinutes, err := app.ParseUpdateReminderMessage(message); err == nil {
		return updateReminderFromMessage(c, phone, referenceId, nMinutes, fromTime)
	}
	return utils.ReminderDetails{}, app.ReminderParseError(fmt.Sprintf("Unable to create reminder from request %s", message))
}

func createReminderFromMessage(c client.Client, phone string, reminderName string, reminderText string, nMinutes int, fromTime time.Time) (utils.ReminderDetails, error) {
	input := utils.ReminderInput{
		FromTime:     fromTime,
		NMinutes:     nMinutes,
		ReminderText: reminderText,
		ReminderName: reminderName,
		Phone:        phone,
	}
	reminderInfo, err := workflows.StartWorkflow(c, &input)
	log.Printf("Creating reminder for Phone %s", input.Phone)
	if err != nil {
		log.Printf("Failed to start workflow: %v", err)
		return reminderInfo, err
	}
	log.Printf("Created reminder for workflowId %s runId %s", reminderInfo.WorkflowId, reminderInfo.RunId)
	err = whatsapp.GetWhatsappClient().SendMessage(
		phone,
		fmt.Sprintf(
			"Scheduled reminder %s: %s to remind at %s. Reference ID=%s",
			reminderInfo.ReminderName,
			reminderInfo.ReminderText,
			utils.GetReminderTime(reminderInfo.FromTime, reminderInfo.NMinutes).Format(app.TIME_FORMAT),
			reminderInfo.ReferenceId,
		),
	)
	return reminderInfo, err
}

func updateReminderFromMessage(c client.Client, phone string, referenceId string, nMinutes int, fromTime time.Time) (utils.ReminderDetails, error) {
	workflowId, runId, err := utils.GetInternalIdsFromReferenceId(referenceId)
	if err != nil {
		log.Printf("Failed to update workflow; unrecognized reference ID: %s", referenceId)
		return utils.ReminderDetails{}, err
	}
	log.Printf("Updating reminder for Phone %s. workflowId=%s runId=%s", phone, workflowId, runId)
	input := utils.ReminderInput{
		FromTime: fromTime,
		NMinutes: nMinutes,
		Phone:    phone,
	}
	reminderDetails, err := workflows.UpdateWorkflow(c, workflowId, runId, &input)
	if err != nil {
		log.Printf("Failed to update workflow: %v", err)
		return utils.ReminderDetails{}, err
	}
	log.Printf("Updated reminder for workflowId %s runId %s", reminderDetails.WorkflowId, reminderDetails.RunId)
	err = whatsapp.GetWhatsappClient().SendMessage(
		phone,
		fmt.Sprintf(
			"Updated reminder %s: %s at %s. referenceId=%s",
			reminderDetails.ReminderName,
			reminderDetails.ReminderText,
			utils.GetReminderTime(reminderDetails.FromTime, reminderDetails.NMinutes).Format(app.TIME_FORMAT),
			reminderDetails.ReferenceId,
		),
	)
	return reminderDetails, err
}

func sendErrorMessage(wc whatsapp.IWhatsappClient, phone string, message string) {
	wc.SendMessage(phone, fmt.Sprintf(
		`Error creating reminder: "%s". Please use the format "New Reminder <Reminder Name>: <Reminder Text>: <1H 30M>"`,
		message,
	))
}

type RequestHandler struct {
	c IWorkflowClient
}

type IWorkflowClient interface {
	GetClient() (client.Client, error)
	client.Client
}

type WorkflowClient struct {
	client.Client
}

func (w WorkflowClient) GetClient() (client.Client, error) {
	return client.NewClient(client.Options{})
}

func (h RequestHandler) HandleList(writer http.ResponseWriter, reader *http.Request) {
	h.DeleteReminderHandler(writer, reader)
}

func (h RequestHandler) HandleCreate(writer http.ResponseWriter, reader *http.Request) {
	h.CreateReminderHandler(writer, reader)
}

func (h RequestHandler) HandleGet(writer http.ResponseWriter, reader *http.Request) {
	h.GetReminderHandler(writer, reader)
}

func (h RequestHandler) HandleUpdate(writer http.ResponseWriter, reader *http.Request) {
	h.UpdateReminderHandler(writer, reader)
}

func (h RequestHandler) HandleDelete(writer http.ResponseWriter, reader *http.Request) {
	h.DeleteReminderHandler(writer, reader)
}

func (h RequestHandler) HandleWhatsappCallback(writer http.ResponseWriter, reader *http.Request) {
	h.WhatsappResponseHandler(writer, reader)
}

func main() {
	r := mux.NewRouter()
	requestHandler := RequestHandler{WorkflowClient{}}
	r.HandleFunc("/reminders", requestHandler.HandleList).Methods("GET")
	r.HandleFunc("/reminders", requestHandler.HandleCreate).Methods("POST")
	r.HandleFunc("/reminders/{referenceId}", requestHandler.HandleGet).Methods("GET")
	r.HandleFunc("/reminders/{referenceId}", requestHandler.HandleUpdate).Methods("PUT")
	r.HandleFunc("/reminders/{referenceId}", requestHandler.HandleDelete).Methods("DELETE")
	r.HandleFunc("/external/reminders/whatsapp", requestHandler.HandleWhatsappCallback).Methods("GET")
	r.HandleFunc("/external/reminders/whatsapp", requestHandler.HandleWhatsappCallback).Methods("POST")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8000", r))
}
