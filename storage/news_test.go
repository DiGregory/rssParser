package storage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var goodTestCases = []*News{
	{
		Title:       "fake title1",
		Description: "fake description1",
		Link:        "yandex.ru",
	},
	{
		Title:       "fake title2",
		Description: "fake description2",
		Link:        "yandex.ru",
	},
	{
		Title:       "fake title2",
		Description: "fake description2",
		Link:        "yandex.ru",
	},
	{
		Title:       "fake title3",
		Description: "fake description3",
		Link:        "yandex.ru",
	},
}

func TestShouldInsertNews(t *testing.T) {
	db, mock, err := sqlmock.New()
	mockedStorage := NewNewsStorage(db)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	for _, t := range goodTestCases {
		mock.ExpectExec("INSERT INTO news\\(title, description, link\\) VALUES \\(\\$1, \\$2, \\$3\\) ON CONFLICT DO NOTHING;").
			WithArgs(t.Title, t.Description, t.Link).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mock.ExpectCommit()

	// now we execute our method
	if err = mockedStorage.CreateNews(goodTestCases); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldSelectNews(t *testing.T) {
	db, mock, err := sqlmock.New()
	mockedStorage := NewNewsStorage(db)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	//args
	limit := int32(2)
	offset := int32(0)

	mock.ExpectQuery("SELECT \\* FROM  news LIMIT coalesce\\(\\$1, 10\\) OFFSET coalesce\\(\\$2, 0\\);").
		WithArgs(limit, offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "link"}).
			AddRow(1, "3", "3", "4"))

	if _, err = mockedStorage.GetNews(&limit, &offset); err != nil {
		t.Errorf("error was not expected while selecting news: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
