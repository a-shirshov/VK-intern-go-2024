package delivery

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"vk-intern_test-case/internal/film/mock"
	"vk-intern_test-case/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func prepareTestEnvironment() *httptest.ResponseRecorder {
	responseRecorder := httptest.NewRecorder()
	return responseRecorder
}

type addFilmTest struct {
	name               string
	inputBodyJSON      string
	beforeTest         func(mockFilmRepository *mock.MockFilmRepository)
	expectedFilmJSON   string
	expectedStatusCode int
}

var addFilmTests = []addFilmTest{
	{
		"Successfully add new Film",
		`{ 
			"title":"Titanic", 
			"description": "Cool film",
			"release_date": "2001-08-06",
			"rating": 8
		}`,
		func(mockFilmRepository *mock.MockFilmRepository) {
			mockFilmRepository.EXPECT().
				AddFilm(&models.FilmWithActors{
					Film: models.Film{
						FilmRequest: models.FilmRequest{
							Title:       "Titanic",
							Description: "Cool film",
							ReleaseDate: "2001-08-06",
							Rating:      8,
						},
					},
				}).
				Return(
					&models.Film{
						ID: 1,
						FilmRequest: models.FilmRequest{
							Title:       "Titanic",
							Description: "Cool film",
							ReleaseDate: "2001-08-06",
							Rating:      8,
						},
					},
					nil)
		},
		`{ 
			"id": 1,
			"title":"Titanic", 
			"description": "Cool film",
			"release_date": "2001-08-06",
			"rating": 8
		}`,
		http.StatusOK,
	},
}

func TestAddFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range addFilmTests {
		t.Run(test.name, func(t *testing.T) {
			mockFilmRepository := mock.NewMockFilmRepository(ctrl)
			filmDeliveryTest := NewFilmDelivery(mockFilmRepository)
			if test.beforeTest != nil {
				test.beforeTest(mockFilmRepository)
			}
			responseRecorder := prepareTestEnvironment()

			request, err := http.NewRequest(http.MethodPost, "/films", strings.NewReader(test.inputBodyJSON))
			assert.Nil(t, err)

			filmDeliveryTest.HandleFilms(responseRecorder, request)
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

type updateFilmTest struct {
	name               string
	queryFilmID        int
	inputBodyJSON      string
	beforeTest         func(mockFilmRepository *mock.MockFilmRepository)
	expectedJSON       string
	expectedStatusCode int
}

var updateFilmTests = []updateFilmTest{
	{
		"Successfully update a Film",
		1,
		`{ 
			"id": 1,
			"title":"Titanic", 
			"description": "Cool film",
			"release_date": "2001-08-06",
			"rating": 8
		}`,
		func(mockFilmRepository *mock.MockFilmRepository) {
			mockFilmRepository.EXPECT().
				UpdateFilm(1, &models.Film{
					ID: 1,
					FilmRequest: models.FilmRequest{
						Title:       "Titanic",
						Description: "Cool film",
						ReleaseDate: "2001-08-06",
						Rating:      8,
					},
				}).
				Return(nil)
		},
		`{
			"status": "OK"
		}`,
		http.StatusOK,
	},
}

func TestUpdateFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range updateFilmTests {
		t.Run(test.name, func(t *testing.T) {
			mockFilmRepository := mock.NewMockFilmRepository(ctrl)
			filmDeliveryTest := NewFilmDelivery(mockFilmRepository)
			if test.beforeTest != nil {
				test.beforeTest(mockFilmRepository)
			}
			responseRecorder := prepareTestEnvironment()

			request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%d", "/films", 1), strings.NewReader(test.inputBodyJSON))
			assert.Nil(t, err)

			filmDeliveryTest.HandleFilms(responseRecorder, request)
			result := responseRecorder.Result()
			defer result.Body.Close()
			data, err := io.ReadAll(result.Body)

			assert.Nil(t, err)
			assert.Equal(t, test.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, "application/json", result.Header[http.CanonicalHeaderKey("content-type")][0])
			assert.JSONEq(t, test.expectedJSON, string(data))
		})
	}
}

type deleteFilmTest struct {
	name               string
	queryFilmID        int
	beforeTest         func(mockFilmRepository *mock.MockFilmRepository)
	expectedJSON       string
	expectedStatusCode int
}

var deleteFilmTests = []deleteFilmTest{
	{
		"Successfully delete a Film",
		1,
		func(mockFilmRepository *mock.MockFilmRepository) {
			mockFilmRepository.EXPECT().
				DeleteFilm(1).
				Return(nil)
		},
		`{
			"status": "OK"
		}`,
		http.StatusOK,
	},
}

func TestDeleteFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range deleteFilmTests {
		t.Run(test.name, func(t *testing.T) {
			mockFilmRepository := mock.NewMockFilmRepository(ctrl)
			filmDeliveryTest := NewFilmDelivery(mockFilmRepository)
			if test.beforeTest != nil {
				test.beforeTest(mockFilmRepository)
			}
			responseRecorder := prepareTestEnvironment()

			request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", "/films", 1), nil)
			assert.Nil(t, err)

			filmDeliveryTest.HandleFilms(responseRecorder, request)
			result := responseRecorder.Result()
			defer result.Body.Close()
			data, err := io.ReadAll(result.Body)

			assert.Nil(t, err)
			assert.Equal(t, test.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, "application/json", result.Header[http.CanonicalHeaderKey("content-type")][0])
			assert.JSONEq(t, test.expectedJSON, string(data))
		})
	}
}

type getFilmsTest struct {
	name               string
	querySortBy        string
	beforeTest         func(mockFilmRepository *mock.MockFilmRepository)
	expectedJSON       string
	expectedStatusCode int
}

var getFilmsTests = []getFilmsTest{
	{
		"Successfully get a list of Film with no query param",
		"",
		func(mockFilmRepository *mock.MockFilmRepository) {
			mockFilmRepository.EXPECT().
				GetFilmsSorted("").
				Return([]models.Film{
					{
						ID: 1,
						FilmRequest: models.FilmRequest{
							Title:       "Titanic",
							Description: "Cool film",
							ReleaseDate: "2001-08-06",
							Rating:      8,
						},
					},
					{
						ID: 2,
						FilmRequest: models.FilmRequest{
							Title:       "Titanic 2",
							Description: "Not Cool film",
							ReleaseDate: "2001-08-06",
							Rating:      7,
						},
					},
				}, nil)
		},
		`[
			{ 
				"id": 1,
				"title":"Titanic", 
				"description": "Cool film",
				"release_date": "2001-08-06",
				"rating": 8
			},
			{ 
				"id": 2,
				"title":"Titanic 2", 
				"description": "Not Cool film",
				"release_date": "2001-08-06",
				"rating": 7
			}
		]`,
		http.StatusOK,
	},
}

func TestGetFilms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range getFilmsTests {
		t.Run(test.name, func(t *testing.T) {
			mockFilmRepository := mock.NewMockFilmRepository(ctrl)
			filmDeliveryTest := NewFilmDelivery(mockFilmRepository)
			if test.beforeTest != nil {
				test.beforeTest(mockFilmRepository)
			}
			responseRecorder := prepareTestEnvironment()

			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?sort_by=%s", "/films", test.querySortBy), nil)
			assert.Nil(t, err)

			filmDeliveryTest.HandleFilms(responseRecorder, request)
			result := responseRecorder.Result()
			defer result.Body.Close()
			data, err := io.ReadAll(result.Body)

			assert.Nil(t, err)
			assert.Equal(t, test.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, "application/json", result.Header[http.CanonicalHeaderKey("content-type")][0])
			assert.JSONEq(t, test.expectedJSON, string(data))
		})
	}
}

type getFilmByTitleTest struct {
	name               string
	queryTitle         string
	beforeTest         func(mockFilmRepository *mock.MockFilmRepository)
	expectedJSON       string
	expectedStatusCode int
}

var getFilmByTitleTests = []getFilmByTitleTest{
	{
		"Successfully get a list of Film by title",
		"Tit",
		func(mockFilmRepository *mock.MockFilmRepository) {
			mockFilmRepository.EXPECT().
				GetFilmsByTitle("Tit").
				Return([]models.Film{
					{
						ID: 1,
						FilmRequest: models.FilmRequest{
							Title:       "Titanic",
							Description: "Cool film",
							ReleaseDate: "2001-08-06",
							Rating:      8,
						},
					},
					{
						ID: 2,
						FilmRequest: models.FilmRequest{
							Title:       "Titanic 2",
							Description: "Not Cool film",
							ReleaseDate: "2001-08-06",
							Rating:      7,
						},
					},
				}, nil)
		},
		`[
			{ 
				"id": 1,
				"title":"Titanic", 
				"description": "Cool film",
				"release_date": "2001-08-06",
				"rating": 8
			},
			{ 
				"id": 2,
				"title":"Titanic 2", 
				"description": "Not Cool film",
				"release_date": "2001-08-06",
				"rating": 7
			}
		]`,
		http.StatusOK,
	},
	{
		"Get an error Film by title",
		"Tit",
		func(mockFilmRepository *mock.MockFilmRepository) {
			mockFilmRepository.EXPECT().
				GetFilmsByTitle("Tit").
				Return([]models.Film{}, errors.New("error text"))
		},
		`{
			"status": "error text"
		}`,
		http.StatusInternalServerError,
	},
}

func TestGetFilmsByTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range getFilmByTitleTests {
		t.Run(test.name, func(t *testing.T) {
			mockFilmRepository := mock.NewMockFilmRepository(ctrl)
			filmDeliveryTest := NewFilmDelivery(mockFilmRepository)
			if test.beforeTest != nil {
				test.beforeTest(mockFilmRepository)
			}
			responseRecorder := prepareTestEnvironment()

			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?title=%s", "/film", test.queryTitle), nil)
			assert.Nil(t, err)

			filmDeliveryTest.HandleFilm(responseRecorder, request)
			result := responseRecorder.Result()
			defer result.Body.Close()
			data, err := io.ReadAll(result.Body)

			assert.Nil(t, err)
			assert.Equal(t, test.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, "application/json", result.Header[http.CanonicalHeaderKey("content-type")][0])
			assert.JSONEq(t, test.expectedJSON, string(data))
		})
	}
}

type getFilmByActorTest struct {
	name               string
	queryActor         string
	beforeTest         func(mockFilmRepository *mock.MockFilmRepository)
	expectedJSON       string
	expectedStatusCode int
}

var getFilmByActorTests = []getFilmByActorTest{
	{
		"Successfully get a list of Film by actor name",
		"Лео",
		func(mockFilmRepository *mock.MockFilmRepository) {
			mockFilmRepository.EXPECT().
				GetFilmsByActor("Лео").
				Return([]models.Film{
					{
						ID: 1,
						FilmRequest: models.FilmRequest{
							Title:       "Titanic",
							Description: "Cool film",
							ReleaseDate: "2001-08-06",
							Rating:      8,
						},
					},
					{
						ID: 2,
						FilmRequest: models.FilmRequest{
							Title:       "Titanic 2",
							Description: "Not Cool film",
							ReleaseDate: "2001-08-06",
							Rating:      7,
						},
					},
				}, nil)
		},
		`[
			{ 
				"id": 1,
				"title":"Titanic", 
				"description": "Cool film",
				"release_date": "2001-08-06",
				"rating": 8
			},
			{ 
				"id": 2,
				"title":"Titanic 2", 
				"description": "Not Cool film",
				"release_date": "2001-08-06",
				"rating": 7
			}
		]`,
		http.StatusOK,
	},
}

func TestGetFilmsByActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range getFilmByActorTests {
		t.Run(test.name, func(t *testing.T) {
			mockFilmRepository := mock.NewMockFilmRepository(ctrl)
			filmDeliveryTest := NewFilmDelivery(mockFilmRepository)
			if test.beforeTest != nil {
				test.beforeTest(mockFilmRepository)
			}
			responseRecorder := prepareTestEnvironment()

			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?actor=%s", "/film", test.queryActor), nil)
			assert.Nil(t, err)

			filmDeliveryTest.HandleFilm(responseRecorder, request)
			result := responseRecorder.Result()
			defer result.Body.Close()
			data, err := io.ReadAll(result.Body)

			assert.Nil(t, err)
			assert.Equal(t, test.expectedStatusCode, responseRecorder.Code)
			assert.Equal(t, "application/json", result.Header[http.CanonicalHeaderKey("content-type")][0])
			assert.JSONEq(t, test.expectedJSON, string(data))
		})
	}
}