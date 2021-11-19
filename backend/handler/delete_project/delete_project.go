package delete_project

import (
	"context"
	"encoding/json"
	"errors"
	"functions/backend/api"
	"functions/backend/config"
	"github.com/hasura/go-graphql-client"
	"github.com/qovery/qovery-client-go"
	"io/ioutil"
	"log"
	"net/http"
)

type ActionPayload struct {
	SessionVariables map[string]string `json:"session_variables"`
	Input            DeleteProjectArgs `json:"input"`
}

type GraphQLError struct {
	Message string `json:"message"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// set the response header as JSON
	w.Header().Set("Content-Type", "application/json")

	// read request body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// parse the body as action payload
	var actionPayload ActionPayload
	err = json.Unmarshal(reqBody, &actionPayload)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// getting user id
	userId := actionPayload.SessionVariables["x-hasura-user-id"]
	if len(userId) == 0 {
		errorObject := GraphQLError{
			Message: "user not authenticated",
		}
		errorBody, _ := json.Marshal(errorObject)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorBody)
		return
	}

	// Send the request params to the Action's generated handler function
	result, err := deleteProject(actionPayload.Input, userId)

	// throw if an error happens
	if err != nil {
		errorObject := GraphQLError{
			Message: err.Error(),
		}
		errorBody, _ := json.Marshal(errorObject)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorBody)
		return
	}

	// Write the response as JSON
	data, _ := json.Marshal(result)
	w.Write(data)
}

func deleteProject(args DeleteProjectArgs, userId string) (response DeleteProjectOutput, err error) {
	log.Printf("received delete project request %v", args)

	response = DeleteProjectOutput{
		Ok: false,
	}

	// try to delete a project using Qovery API
	err = callQoveryApi(args.Input.Id)
	if err != nil {
		return response, err
	}

	response.Ok = true

	return response, nil
}

func callQoveryApi(id int32) error {
	cfg := qovery.NewConfiguration()
	cfg.AddDefaultHeader("Authorization", "Bearer "+config.QoveryApiToken)
	client := qovery.NewAPIClient(cfg)

	var query struct {
		Project []struct {
			Id                graphql.Int
			Qovery_Project_Id graphql.String `json:"qovery_project_id"`
		} `graphql:"project(where: {id: {_eq: $id}})"`
	}
	vars := map[string]interface{}{
		"id": graphql.Int(id),
	}

	// getting project by id
	err := api.HasuraClient.Query(context.Background(), &query, vars)
	if err != nil {
		return err
	}
	if len(query.Project) != 1 {
		return errors.New("project not found")
	}

	res, err := client.ProjectMainCallsApi.DeleteProject(context.Background(), string(query.Project[0].Qovery_Project_Id)).Execute()
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return errors.New("received " + res.Status + " deleting a project from Qovery API")
	}

	var mutation struct {
		DeleteProjectByPk struct {
			Id graphql.Int
		} `graphql:"delete_project_by_pk(id: $id)"`
	}
	vars = map[string]interface{}{
		"id": graphql.Int(id),
	}

	// trying to create a new project in Hasura backend
	err = api.HasuraClient.Mutate(context.Background(), &mutation, vars)
	if err != nil {
		return err
	}

	return nil
}
