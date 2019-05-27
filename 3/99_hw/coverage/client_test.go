package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestSearchServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	cases := []struct {
		expected  *SearchResponse
		SearchReq SearchRequest
		Err       bool
	}{
		{
			SearchReq: SearchRequest{
				Limit:      -1,
				Offset:     0,
				Query:      "am",
				OrderField: "age",
				OrderBy:    OrderByDesc,
			},
			Err: true,
		},
		{
			SearchReq: SearchRequest{
				Limit:      0,
				Offset:     0,
				Query:      "i",
				OrderField: "age",
				OrderBy:    OrderByAsc,
			},
			Err: true,
		},
		{
			SearchReq: SearchRequest{
				Limit:      2,
				Offset:     1,
				Query:      "i",
				OrderField: "qw",
				OrderBy:    OrderByAsc,
			},
			Err: true,
		},
		{
			SearchReq: SearchRequest{
				Limit:      3,
				Offset:     4,
				Query:      "anim a",
				OrderField: "id",
				OrderBy:    OrderByAsc,
			},
			Err: true,
		},
		{
			SearchReq: SearchRequest{
				Limit:      2,
				Offset:     1,
				Query:      "i",
				OrderField: "age",
				OrderBy:    -5,
			},
			Err: true,
		},
		{
			expected: &SearchResponse{
				[]User{
					{
						Id:     11,
						Age:    32,
						Name:   "Gilmore Guerra",
						About:  "Labore consectetur do sit et mollit non incididunt. Amet aute voluptate enim et sit Lorem elit. Fugiat proident ullamco ullamco sint pariatur deserunt eu nulla consectetur culpa eiusmod. Veniam irure et deserunt consectetur incididunt ad ipsum sint. Consectetur voluptate adipisicing aute fugiat aliquip culpa qui nisi ut ex esse ex. Sint et anim aliqua pariatur.\n",
						Gender: "male",
					},
					{
						Id:     12,
						Age:    36,
						Name:   "Cruz Guerrero",
						About:  "Sunt enim ad fugiat minim id esse proident laborum magna magna. Velit anim aliqua nulla laborum consequat veniam reprehenderit enim fugiat ipsum mollit nisi. Nisi do reprehenderit aute sint sit culpa id Lorem proident id tempor. Irure ut ipsum sit non quis aliqua in voluptate magna. Ipsum non aliquip quis incididunt incididunt aute sint. Minim dolor in mollit aute duis consectetur.\n",
						Gender: "male",
					},
					{
						Id:     25,
						Age:    32,
						Name:   "Katheryn Jacobs",
						About:  "Magna excepteur anim amet id consequat tempor dolor sunt id enim ipsum ea est ex. In do ea sint qui in minim mollit anim est et minim dolore velit laborum. Officia commodo duis ut proident laboris fugiat commodo do ex duis consequat exercitation. Ad et excepteur ex ea exercitation id fugiat exercitation amet proident adipisicing laboris id deserunt. Commodo proident laborum elit ex aliqua labore culpa ullamco occaecat voluptate voluptate laboris deserunt magna.\n",
						Gender: "female",
					},
				},
				false,
			},

			SearchReq: SearchRequest{
				Limit:      30,
				Offset:     0,
				Query:      "anim a",
				OrderField: "id",
				OrderBy:    OrderByAsc,
			},
			Err: false,
		},
		{
			expected: &SearchResponse{
				[]User{
					{
						Id:     11,
						Age:    32,
						Name:   "Gilmore Guerra",
						About:  "Labore consectetur do sit et mollit non incididunt. Amet aute voluptate enim et sit Lorem elit. Fugiat proident ullamco ullamco sint pariatur deserunt eu nulla consectetur culpa eiusmod. Veniam irure et deserunt consectetur incididunt ad ipsum sint. Consectetur voluptate adipisicing aute fugiat aliquip culpa qui nisi ut ex esse ex. Sint et anim aliqua pariatur.\n",
						Gender: "male",
					},
					{
						Id:     12,
						Age:    36,
						Name:   "Cruz Guerrero",
						About:  "Sunt enim ad fugiat minim id esse proident laborum magna magna. Velit anim aliqua nulla laborum consequat veniam reprehenderit enim fugiat ipsum mollit nisi. Nisi do reprehenderit aute sint sit culpa id Lorem proident id tempor. Irure ut ipsum sit non quis aliqua in voluptate magna. Ipsum non aliquip quis incididunt incididunt aute sint. Minim dolor in mollit aute duis consectetur.\n",
						Gender: "male",
					},
					{
						Id:     25,
						Age:    32,
						Name:   "Katheryn Jacobs",
						About:  "Magna excepteur anim amet id consequat tempor dolor sunt id enim ipsum ea est ex. In do ea sint qui in minim mollit anim est et minim dolore velit laborum. Officia commodo duis ut proident laboris fugiat commodo do ex duis consequat exercitation. Ad et excepteur ex ea exercitation id fugiat exercitation amet proident adipisicing laboris id deserunt. Commodo proident laborum elit ex aliqua labore culpa ullamco occaecat voluptate voluptate laboris deserunt magna.\n",
						Gender: "female",
					},
				},
				false,
			},

			SearchReq: SearchRequest{
				Limit:      30,
				Offset:     0,
				Query:      "anim a",
				OrderField: "id",
				OrderBy:    OrderByAsIs,
			},
			Err: false,
		},
		{
			expected: &SearchResponse{
				[]User{
					{
						Id:     25,
						Age:    32,
						Name:   "Katheryn Jacobs",
						About:  "Magna excepteur anim amet id consequat tempor dolor sunt id enim ipsum ea est ex. In do ea sint qui in minim mollit anim est et minim dolore velit laborum. Officia commodo duis ut proident laboris fugiat commodo do ex duis consequat exercitation. Ad et excepteur ex ea exercitation id fugiat exercitation amet proident adipisicing laboris id deserunt. Commodo proident laborum elit ex aliqua labore culpa ullamco occaecat voluptate voluptate laboris deserunt magna.\n",
						Gender: "female",
					},
					{
						Id:     12,
						Age:    36,
						Name:   "Cruz Guerrero",
						About:  "Sunt enim ad fugiat minim id esse proident laborum magna magna. Velit anim aliqua nulla laborum consequat veniam reprehenderit enim fugiat ipsum mollit nisi. Nisi do reprehenderit aute sint sit culpa id Lorem proident id tempor. Irure ut ipsum sit non quis aliqua in voluptate magna. Ipsum non aliquip quis incididunt incididunt aute sint. Minim dolor in mollit aute duis consectetur.\n",
						Gender: "male",
					},
					{
						Id:     11,
						Age:    32,
						Name:   "Gilmore Guerra",
						About:  "Labore consectetur do sit et mollit non incididunt. Amet aute voluptate enim et sit Lorem elit. Fugiat proident ullamco ullamco sint pariatur deserunt eu nulla consectetur culpa eiusmod. Veniam irure et deserunt consectetur incididunt ad ipsum sint. Consectetur voluptate adipisicing aute fugiat aliquip culpa qui nisi ut ex esse ex. Sint et anim aliqua pariatur.\n",
						Gender: "male",
					},
				},
				false,
			},

			SearchReq: SearchRequest{
				Limit:      30,
				Offset:     0,
				Query:      "anim a",
				OrderField: "id",
				OrderBy:    OrderByDesc,
			},
			Err: false,
		},
		{
			expected: &SearchResponse{
				[]User{
					{
						Id:     25,
						Age:    32,
						Name:   "Katheryn Jacobs",
						About:  "Magna excepteur anim amet id consequat tempor dolor sunt id enim ipsum ea est ex. In do ea sint qui in minim mollit anim est et minim dolore velit laborum. Officia commodo duis ut proident laboris fugiat commodo do ex duis consequat exercitation. Ad et excepteur ex ea exercitation id fugiat exercitation amet proident adipisicing laboris id deserunt. Commodo proident laborum elit ex aliqua labore culpa ullamco occaecat voluptate voluptate laboris deserunt magna.\n",
						Gender: "female",
					},
					{
						Id:     11,
						Age:    32,
						Name:   "Gilmore Guerra",
						About:  "Labore consectetur do sit et mollit non incididunt. Amet aute voluptate enim et sit Lorem elit. Fugiat proident ullamco ullamco sint pariatur deserunt eu nulla consectetur culpa eiusmod. Veniam irure et deserunt consectetur incididunt ad ipsum sint. Consectetur voluptate adipisicing aute fugiat aliquip culpa qui nisi ut ex esse ex. Sint et anim aliqua pariatur.\n",
						Gender: "male",
					},
					{
						Id:     12,
						Age:    36,
						Name:   "Cruz Guerrero",
						About:  "Sunt enim ad fugiat minim id esse proident laborum magna magna. Velit anim aliqua nulla laborum consequat veniam reprehenderit enim fugiat ipsum mollit nisi. Nisi do reprehenderit aute sint sit culpa id Lorem proident id tempor. Irure ut ipsum sit non quis aliqua in voluptate magna. Ipsum non aliquip quis incididunt incididunt aute sint. Minim dolor in mollit aute duis consectetur.\n",
						Gender: "male",
					},
				},
				false,
			},

			SearchReq: SearchRequest{
				Limit:      3,
				Offset:     0,
				Query:      "anim a",
				OrderField: "",
				OrderBy:    OrderByDesc,
			},
			Err: false,
		},
		{
			expected: &SearchResponse{
				[]User{
					{
						Id:     12,
						Age:    36,
						Name:   "Cruz Guerrero",
						About:  "Sunt enim ad fugiat minim id esse proident laborum magna magna. Velit anim aliqua nulla laborum consequat veniam reprehenderit enim fugiat ipsum mollit nisi. Nisi do reprehenderit aute sint sit culpa id Lorem proident id tempor. Irure ut ipsum sit non quis aliqua in voluptate magna. Ipsum non aliquip quis incididunt incididunt aute sint. Minim dolor in mollit aute duis consectetur.\n",
						Gender: "male",
					},
					{
						Id:     11,
						Age:    32,
						Name:   "Gilmore Guerra",
						About:  "Labore consectetur do sit et mollit non incididunt. Amet aute voluptate enim et sit Lorem elit. Fugiat proident ullamco ullamco sint pariatur deserunt eu nulla consectetur culpa eiusmod. Veniam irure et deserunt consectetur incididunt ad ipsum sint. Consectetur voluptate adipisicing aute fugiat aliquip culpa qui nisi ut ex esse ex. Sint et anim aliqua pariatur.\n",
						Gender: "male",
					},
					{
						Id:     25,
						Age:    32,
						Name:   "Katheryn Jacobs",
						About:  "Magna excepteur anim amet id consequat tempor dolor sunt id enim ipsum ea est ex. In do ea sint qui in minim mollit anim est et minim dolore velit laborum. Officia commodo duis ut proident laboris fugiat commodo do ex duis consequat exercitation. Ad et excepteur ex ea exercitation id fugiat exercitation amet proident adipisicing laboris id deserunt. Commodo proident laborum elit ex aliqua labore culpa ullamco occaecat voluptate voluptate laboris deserunt magna.\n",
						Gender: "female",
					},
				},
				false,
			},

			SearchReq: SearchRequest{
				Limit:      3,
				Offset:     0,
				Query:      "anim a",
				OrderField: "name",
				OrderBy:    OrderByAsc,
			},
			Err: false,
		},
		{
			SearchReq: SearchRequest{
				Limit:      2,
				Offset:     -1,
				Query:      "i",
				OrderField: "age",
				OrderBy:    OrderByAsc,
			},
			Err: true,
		},
		{
			expected: &SearchResponse{
				[]User{
					{
						Id:     12,
						Age:    36,
						Name:   "Cruz Guerrero",
						About:  "Sunt enim ad fugiat minim id esse proident laborum magna magna. Velit anim aliqua nulla laborum consequat veniam reprehenderit enim fugiat ipsum mollit nisi. Nisi do reprehenderit aute sint sit culpa id Lorem proident id tempor. Irure ut ipsum sit non quis aliqua in voluptate magna. Ipsum non aliquip quis incididunt incididunt aute sint. Minim dolor in mollit aute duis consectetur.\n",
						Gender: "male",
					},
					{
						Id:     11,
						Age:    32,
						Name:   "Gilmore Guerra",
						About:  "Labore consectetur do sit et mollit non incididunt. Amet aute voluptate enim et sit Lorem elit. Fugiat proident ullamco ullamco sint pariatur deserunt eu nulla consectetur culpa eiusmod. Veniam irure et deserunt consectetur incididunt ad ipsum sint. Consectetur voluptate adipisicing aute fugiat aliquip culpa qui nisi ut ex esse ex. Sint et anim aliqua pariatur.\n",
						Gender: "male",
					},
				},
				true,
			},

			SearchReq: SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "anim a",
				OrderField: "age",
				OrderBy:    OrderByDesc,
			},
			Err: false,
		},
		{
			expected: &SearchResponse{
				[]User{
					{
						Id:     11,
						Age:    32,
						Name:   "Gilmore Guerra",
						About:  "Labore consectetur do sit et mollit non incididunt. Amet aute voluptate enim et sit Lorem elit. Fugiat proident ullamco ullamco sint pariatur deserunt eu nulla consectetur culpa eiusmod. Veniam irure et deserunt consectetur incididunt ad ipsum sint. Consectetur voluptate adipisicing aute fugiat aliquip culpa qui nisi ut ex esse ex. Sint et anim aliqua pariatur.\n",
						Gender: "male",
					},
					{
						Id:     25,
						Age:    32,
						Name:   "Katheryn Jacobs",
						About:  "Magna excepteur anim amet id consequat tempor dolor sunt id enim ipsum ea est ex. In do ea sint qui in minim mollit anim est et minim dolore velit laborum. Officia commodo duis ut proident laboris fugiat commodo do ex duis consequat exercitation. Ad et excepteur ex ea exercitation id fugiat exercitation amet proident adipisicing laboris id deserunt. Commodo proident laborum elit ex aliqua labore culpa ullamco occaecat voluptate voluptate laboris deserunt magna.\n",
						Gender: "female",
					},
					{
						Id:     12,
						Age:    36,
						Name:   "Cruz Guerrero",
						About:  "Sunt enim ad fugiat minim id esse proident laborum magna magna. Velit anim aliqua nulla laborum consequat veniam reprehenderit enim fugiat ipsum mollit nisi. Nisi do reprehenderit aute sint sit culpa id Lorem proident id tempor. Irure ut ipsum sit non quis aliqua in voluptate magna. Ipsum non aliquip quis incididunt incididunt aute sint. Minim dolor in mollit aute duis consectetur.\n",
						Gender: "male",
					},
				},
				false,
			},

			SearchReq: SearchRequest{
				Limit:      3,
				Offset:     0,
				Query:      "anim a",
				OrderField: "age",
				OrderBy:    OrderByAsc,
			},
			Err: false,
		},
		{
			SearchReq: SearchRequest{
				Limit:      2,
				Offset:     -1,
				Query:      "i",
				OrderField: "age",
				OrderBy:    OrderByAsc,
			},
			Err: true,
		},
	}
	client := SearchClient{
		AccessToken: Token,
		URL:         ts.URL,
	}
	for caseNum, item := range cases {
		resp, err := client.FindUsers(item.SearchReq)
		if err != nil && !item.Err {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && item.Err {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !reflect.DeepEqual(item.expected, resp) {
			t.Errorf("[%d] wrong result, expected %#v\n, got %#v", caseNum, item.expected, resp)
		}
	}
	ts.Close()
}

func TestSearchClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	cases := []struct {
		SearchCl SearchClient
		Err      bool
	}{
		{
			SearchCl: SearchClient{
				AccessToken: Token,
				URL:         ts.URL,
			},
			Err: false,
		},
		{
			SearchCl: SearchClient{
				AccessToken: Token,
				URL:         "",
			},
			Err: true,
		},
		{
			SearchCl: SearchClient{
				AccessToken: "qwee",
				URL:         ts.URL,
			},
			Err: true,
		},
	}
	req := SearchRequest{
		Limit:      1,
		Offset:     0,
		Query:      "amely",
		OrderField: "age",
		OrderBy:    OrderByDesc,
	}
	for caseNum, item := range cases {
		_, err := item.SearchCl.FindUsers(req)
		if err != nil && !item.Err {
			t.Errorf("[%d] Exptected error %#v", caseNum, err)
		}
	}
	ts.Close()
}

func TestFindFile(t *testing.T) {
	FileName = "wfe.xml"
	defer func() { FileName = "dataset.xml" }()
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	client := SearchClient{
		AccessToken: Token,
		URL:         ts.URL,
	}
	req := SearchRequest{
		Limit:      5,
		Offset:     0,
		Query:      "a",
		OrderField: "age",
		OrderBy:    OrderByDesc,
	}
	_, err := client.FindUsers(req)
	if err == nil {
		t.Errorf("Exptected error")
	}
	ts.Close()
}

func TestClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	client.Timeout = time.Nanosecond
	defer func() { client.Timeout = time.Second }()
	req := SearchRequest{
		Limit:      5,
		Offset:     0,
		Query:      "a",
		OrderField: "age",
		OrderBy:    OrderByDesc,
	}
	client := SearchClient{
		AccessToken: Token,
		URL:         ts.URL,
	}
	_, err := client.FindUsers(req)
	if err == nil {
		t.Errorf("Exptected error")
	}
	ts.Close()
}
