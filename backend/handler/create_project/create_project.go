package create_project

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
	"strconv"
)

type ActionPayload struct {
	SessionVariables map[string]string `json:"session_variables"`
	Input            CreateProjectArgs `json:"input"`
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
	result, err := createProject(actionPayload.Input, userId)

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

func createProject(args CreateProjectArgs, userId string) (response CreateProjectOutput, err error) {
	log.Printf("received create project request %v", args)

	response = CreateProjectOutput{
		Id:   0,
		Name: "",
		Url:  "",
	}

	// try to create a new project using Qovery API
	id, url, err := callQoveryApi(args.Input.Name, userId)
	if err != nil {
		return response, err
	}

	response.Id = int32(id)
	response.Name = args.Input.Name
	response.Url = string(url)

	return response, nil
}

func callQoveryApi(name string, userId string) (graphql.Int, graphql.String, error) {
	cfg := qovery.NewConfiguration()
	cfg.AddDefaultHeader("Authorization", "Bearer "+config.QoveryApiToken)
	client := qovery.NewAPIClient(cfg)

	qp, res, err := client.ProjectsApi.CreateProject(context.Background(), config.QoveryOrganizationId).ProjectRequest(qovery.ProjectRequest{
		Name: name,
	}).Execute()
	if err != nil {
		return 0, "", err
	}
	if res.StatusCode >= 400 {
		return 0, "", errors.New("received " + res.Status + " creating a new project from Qovery API")
	}

	qe, res, err := client.EnvironmentsApi.CreateEnvironment(context.Background(), qp.Id).EnvironmentRequest(qovery.EnvironmentRequest{
		Name: "production",
	}).Execute()
	if err != nil {
		return 0, "", err
	}
	if res.StatusCode >= 400 {
		return 0, "", errors.New("received " + res.Status + " creating a new environment from Qovery API")
	}

	cpu := float32(1000)
	memory := float32(2048)
	buildMode := "DOCKER"
	dockerfilePath := "Dockerfile"
	branch := "main"
	protocol := "HTTP"
	accessible := true
	externalPort := int32(443)
	ports := []qovery.ApplicationPortRequestPorts{
		{
			InternalPort:       1337,
			ExternalPort:       &externalPort,
			PubliclyAccessible: &accessible,
			Protocol:           protocol,
		},
	}
	qa, res, err := client.ApplicationsApi.CreateApplication(context.Background(), qe.Id).ApplicationRequest(qovery.ApplicationRequest{
		Name: "strapi",
		GitRepository: qovery.ApplicationGitRepositoryRequest{
			Url:      "https://github.com/Qovery/strapi.git",
			Branch:   &branch,
			RootPath: "/",
		},
		BuildMode:      &buildMode,
		DockerfilePath: &dockerfilePath,
		Cpu:            &cpu,
		Memory:         &memory,
		Ports:          &ports,
	}).Execute()
	if err != nil {
		return 0, "", err
	}
	if res.StatusCode >= 400 {
		return 0, "", errors.New("received " + res.Status + " creating a new application from Qovery API")
	}

	links, res, err := client.ApplicationMainCallsApi.ListApplicationLinks(context.Background(), qa.Id).Execute()
	if err != nil {
		return 0, "", err
	}
	if res.StatusCode >= 400 {
		return 0, "", errors.New("received " + res.Status + " getting a new application link from Qovery API")
	}
	// TODO create Redis
	// TODO create MariaDB

	var mutation struct {
		InsertProjectOne struct {
			Id graphql.Int
		} `graphql:"insert_project_one(object: {owner_id: $owner_id, qovery_environment_id: $qovery_environment_id, qovery_project_id: $qovery_project_id, url: $url, name: $name})"`
	}
	i, err := strconv.ParseInt(userId, 10, 32)
	ownerId := int32(i)
	vars := map[string]interface{}{
		"owner_id":              graphql.Int(ownerId),
		"qovery_environment_id": graphql.String(qe.Id),
		"qovery_project_id":     graphql.String(qp.Id),
		"url":                   graphql.String(links.GetResults()[0].GetUrl()),
		"name":                  graphql.String(name),
	}

	// trying to create a new project in Hasura backend
	err = api.HasuraClient.Mutate(context.Background(), &mutation, vars)
	if err != nil {
		return 0, "", err
	}

	return mutation.InsertProjectOne.Id, graphql.String(links.GetResults()[0].GetUrl()), nil
}
