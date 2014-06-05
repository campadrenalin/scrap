package scrap

import (
	"bytes"
	"testing"
)

func dummyRetriever(ScraperRequest) (Node, error) {
	return Node{}, nil
}

type sc_valid_test struct {
	Config   ScraperConfig
	ErrorMsg string
}

func (svt sc_valid_test) Run(t *testing.T) {
	err := svt.Config.Validate()
	if svt.ErrorMsg == "" {
		// Should not result in an error
		if err != nil {
			t.Fatal(err)
		}
	} else {
		// Should result in specific error
		if err == nil {
			t.Fatalf("Should have failed (%s), but didn't", svt.ErrorMsg)
		}
		compare(t, svt.ErrorMsg, err.Error())
	}
}

func TestScraperConfig_Validate(t *testing.T) {
	var remarks, debug bytes.Buffer
	tests := []sc_valid_test{
		sc_valid_test{
			ScraperConfig{
				Retriever: dummyRetriever,
				Remarks:   &remarks,
				Debug:     &debug,
			},
			"",
		},
		sc_valid_test{
			ScraperConfig{
				Retriever: nil,
				Remarks:   &remarks,
				Debug:     &debug,
			},
			"ScraperConfig not valid if Retriever == nil",
		},
		sc_valid_test{
			ScraperConfig{
				Retriever: dummyRetriever,
				Remarks:   nil,
				Debug:     &debug,
			},
			"ScraperConfig not valid if Remarks == nil",
		},
		sc_valid_test{
			ScraperConfig{
				Retriever: dummyRetriever,
				Remarks:   &remarks,
				Debug:     nil,
			},
			"ScraperConfig not valid if Debug == nil",
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestScraper_BadConfig(t *testing.T) {
	config := ScraperConfig{}
	_, err := NewScraper(config)
	if err == nil {
		t.Fatal("Should have failed, didn't")
	}
}

func TestScraper_CreateRequest(t *testing.T) {
	var remarks, debug bytes.Buffer
	config := ScraperConfig{
		Retriever: dummyRetriever,
		Remarks:   &remarks,
		Debug:     &debug,
	}
	s, err := NewScraper(config)
	if err != nil {
		t.Fatal(err)
	}
	req := s.CreateRequest("/")

	if req.RequestQueue.(*Scraper) != &s {
		t.Fatal("Should have set req.RequestQueue")
	}
	if req.Url != "/" {
		t.Fatal("Should have set req.Url")
	}

	req.Remarks.Println("This is a remark")
	req.Debug.Println("This is a debug note")

	compare(t, "/: This is a remark\n", remarks.String())
	compare(t, "/: This is a debug note\n", debug.String())
}
