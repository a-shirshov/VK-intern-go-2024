package delivery

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"vk-intern_test-case/internal/actor/mock"
	"vk-intern_test-case/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func prepareTestEnvironment() *httptest.ResponseRecorder {
	responseRecorder := httptest.NewRecorder()
	return responseRecorder
}

type addActorTest struct {
	name               string
	inputBodyJSON      string
	beforeTest         func(mockFilmRepository *mock.MockActorRepository)
	expectedFilmJSON   string
	expectedStatusCode int
}

var addActorTests = []addActorTest{
	{
		"Successfully add new Actor",
		`{ 
			"name":"Леонардо Ди Каприо", 
			"gender": "Мужской",
			"date_of_birth": "2002-07-13"
		}`,
		func(mockActorRepository *mock.MockActorRepository) {
			mockActorRepository.EXPECT().
				AddActor(&models.Actor{
					ActorRequest: models.ActorRequest{
						Name: "Леонардо Ди Каприо",
						Gender: "Мужской",
						DateOfBirth: "2002-07-13",
					},
				}).
				Return(
					&models.Actor{
						ID: 1,
						ActorRequest: models.ActorRequest{
							Name: "Леонардо Ди Каприо",
							Gender: "Мужской",
							DateOfBirth: "2002-07-13",
						},
					},
					nil)
		},
		`{ 
			"id":1,
			"name":"Леонардо Ди Каприо", 
			"gender": "Мужской",
			"date_of_birth": "2002-07-13"
		}`,
		http.StatusOK,
	},
}

func TestAddActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range addActorTests {
		t.Run(test.name, func(t *testing.T) {
			mockActorRepository := mock.NewMockActorRepository(ctrl)
			filmDeliveryTest := NewActorDelivery(mockActorRepository)
			if test.beforeTest != nil {
				test.beforeTest(mockActorRepository)
			}
			responseRecorder := prepareTestEnvironment()

			request, err := http.NewRequest(http.MethodPost, "/actors", strings.NewReader(test.inputBodyJSON))
			assert.Nil(t, err)

			filmDeliveryTest.HandleActors(responseRecorder, request)
			result := responseRecorder.Result()
			defer result.Body.Close()
			data, err := io.ReadAll(result.Body)

			assert.Nil(t, err)
			assert.Equal(t, test.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, "application/json", result.Header[http.CanonicalHeaderKey("content-type")][0])
			assert.JSONEq(t, test.expectedFilmJSON, string(data))
		})
	}
}

type updateActorTest struct {
	name               string
	queryActorID       int
	inputBodyJSON      string
	beforeTest         func(mockFilmRepository *mock.MockActorRepository)
	expectedFilmJSON   string
	expectedStatusCode int
}

var updateActorTests = []updateActorTest{
	{
		"Successfully add new Actor",
		1,
		`{ 
			"name":"Леонардо Ди Каприо", 
			"gender": "Мужской",
			"date_of_birth": "2002-07-13"
		}`,
		func(mockActorRepository *mock.MockActorRepository) {
			mockActorRepository.EXPECT().
				UpdateActor(1, &models.Actor{
					ActorRequest: models.ActorRequest{
						Name: "Леонардо Ди Каприо",
						Gender: "Мужской",
						DateOfBirth: "2002-07-13",
					},
				}).
				Return(nil)
		},
		`{ 
			"status":"OK"
		}`,
		http.StatusOK,
	},
	{
		"Successfully add new Actor",
		1,
		`{ 
			"name":"Леонардо Ди Каприо", 
			"gender": "Мужской",
			"date_of_birth": "2002-07-13"
		}`,
		func(mockActorRepository *mock.MockActorRepository) {
			mockActorRepository.EXPECT().
				UpdateActor(1, &models.Actor{
					ActorRequest: models.ActorRequest{
						Name: "Леонардо Ди Каприо",
						Gender: "Мужской",
						DateOfBirth: "2002-07-13",
					},
				}).
				Return(errors.New("new error"))
		},
		`{ 
			"status":"new error"
		}`,
		http.StatusInternalServerError,
	},
}

func TestUpdateActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range updateActorTests {
		t.Run(test.name, func(t *testing.T) {
			mockActorRepository := mock.NewMockActorRepository(ctrl)
			filmDeliveryTest := NewActorDelivery(mockActorRepository)
			if test.beforeTest != nil {
				test.beforeTest(mockActorRepository)
			}
			responseRecorder := prepareTestEnvironment()

			request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%d", "/actors", test.queryActorID), strings.NewReader(test.inputBodyJSON))
			assert.Nil(t, err)

			filmDeliveryTest.HandleActors(responseRecorder, request)
			result := responseRecorder.Result()
			defer result.Body.Close()
			data, err := io.ReadAll(result.Body)

			assert.Nil(t, err)
			assert.Equal(t, test.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, "application/json", result.Header[http.CanonicalHeaderKey("content-type")][0])
			assert.JSONEq(t, test.expectedFilmJSON, string(data))
		})
	}
}

type deleteActorTest struct {
	name               string
	queryActorID       int
	beforeTest         func(mockFilmRepository *mock.MockActorRepository)
	expectedFilmJSON   string
	expectedStatusCode int
}

var deleteActorTests = []deleteActorTest{
	{
		"Successfully add new Actor",
		1,
		func(mockActorRepository *mock.MockActorRepository) {
			mockActorRepository.EXPECT().
				DeleteActor(1).
				Return(nil)
		},
		`{ 
			"status":"OK"
		}`,
		http.StatusOK,
	},
}

func TestDeleteActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range deleteActorTests {
		t.Run(test.name, func(t *testing.T) {
			mockActorRepository := mock.NewMockActorRepository(ctrl)
			filmDeliveryTest := NewActorDelivery(mockActorRepository)
			if test.beforeTest != nil {
				test.beforeTest(mockActorRepository)
			}
			responseRecorder := prepareTestEnvironment()

			request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", "/actors", test.queryActorID), nil)
			assert.Nil(t, err)

			filmDeliveryTest.HandleActors(responseRecorder, request)
			result := responseRecorder.Result()
			defer result.Body.Close()
			data, err := io.ReadAll(result.Body)

			assert.Nil(t, err)
			assert.Equal(t, test.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, "application/json", result.Header[http.CanonicalHeaderKey("content-type")][0])
			assert.JSONEq(t, test.expectedFilmJSON, string(data))
		})
	}
}

type getActorsTest struct {
	name               string
	beforeTest         func(mockFilmRepository *mock.MockActorRepository)
	expectedFilmJSON   string
	expectedStatusCode int
}

var getActorsTests = []getActorsTest{
	{
		"Successfully get Actors with Films",
		func(mockActorRepository *mock.MockActorRepository) {
			mockActorRepository.EXPECT().
				GetActors().
				Return([]models.ActorWithFilms{
					{
						Actor: models.Actor{
							ID: 1,
							ActorRequest: models.ActorRequest{
								Name: "Леонардо Ди Каприо", 
								Gender: "Мужской", 
								DateOfBirth: "2014-03-18",
							},
						},
						Films: []models.Film{
							{
								ID: 1,
								FilmRequest: models.FilmRequest{
									Title: "Титаник", 
									Description: "cool film", 
									ReleaseDate: "2020-06-10", 
									Rating: 8,
								},
							},
							{
								ID: 2,
								FilmRequest: models.FilmRequest{
									Title: "Титаник 2", 
									Description: "cool film", 
									ReleaseDate: "2022-06-10", 
									Rating: 7,
								},
							},
						},
					},
					{
						Actor: models.Actor{
							ID: 2,
							ActorRequest: models.ActorRequest{
								Name: "Марго Робби", 
								Gender: "Женский", 
								DateOfBirth: "2014-03-18",
							},
						},
						Films: []models.Film{
							{
								ID: 3,
								FilmRequest: models.FilmRequest{
									Title: "Барби", 
									Description: "cool film", 
									ReleaseDate: "2020-06-10", 
									Rating: 8,
								},
							},
						},
					},
				}, nil)
		},
		`[
			{
			  "id": 1,
			  "name": "Леонардо Ди Каприо",
			  "gender": "Мужской",
			  "date_of_birth": "2014-03-18",
			  "films": [
				{
				  "id": 1,
				  "title": "Титаник",
				  "description": "cool film",
				  "release_date": "2020-06-10",
				  "rating": 8
				},
				{
				  "id": 2,
				  "title": "Титаник 2",
				  "description": "cool film",
				  "release_date": "2022-06-10",
				  "rating": 7
				}
			  ]
			},
			{
				"id": 2,
				"name": "Марго Робби",
				"gender": "Женский",
				"date_of_birth": "2014-03-18",
				"films": [
				  {
					"id": 3,
					"title": "Барби",
					"description": "cool film",
					"release_date": "2020-06-10",
					"rating": 8
				  }
				]
			  }
		  ]`,
		http.StatusOK,
	},
}

func TestGetActors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range getActorsTests {
		t.Run(test.name, func(t *testing.T) {
			mockActorRepository := mock.NewMockActorRepository(ctrl)
			filmDeliveryTest := NewActorDelivery(mockActorRepository)
			if test.beforeTest != nil {
				test.beforeTest(mockActorRepository)
			}
			responseRecorder := prepareTestEnvironment()

			request, err := http.NewRequest(http.MethodGet, "/actors", nil)
			assert.Nil(t, err)

			filmDeliveryTest.HandleActors(responseRecorder, request)
			result := responseRecorder.Result()
			defer result.Body.Close()
			data, err := io.ReadAll(result.Body)

			assert.Nil(t, err)
			assert.Equal(t, test.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, "application/json", result.Header[http.CanonicalHeaderKey("content-type")][0])
			assert.JSONEq(t, test.expectedFilmJSON, string(data))
		})
	}
}