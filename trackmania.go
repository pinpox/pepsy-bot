package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/olekukonko/tablewriter"
	"strings"
	"time"
)

func MakeSeasonTable() string {

	seasonInfo, err := getCampaignInfo("60530/45849")
	if err != nil {
		log.Println(err)
		return ""
	}

	leads, err := getCampaignLeaderboard(seasonInfo.Leaderboarduid)

	if err != nil {
		log.Println(err)
		return ""
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Position", "Player", "Points"})

	fmt.Println("CAMPAIGN LEADS")
	for _, v := range leads.Tops {

	    table.Append([]string{
			fmt.Sprintf("%v", v.Position),
			v.Player.Name,
			fmt.Sprintf("%v", v.Points),
		})
	}

	table.Render()
	return fmt.Sprintf("```\n%s\n```", tableString.String())

}


func getMapPack(packID string) (MapPackMapInfo, error) {
	url := fmt.Sprintf("https://trackmania.exchange/api/mappack/get_mappack_tracks/%s", packID)
	var mpack MapPackMapInfo
	return getTMObject(url, mpack)
}

func getMapLeaderboard(mapID string) (Leaderboard, error) {
	UrlLeaderboards := fmt.Sprintf("https://trackmania.io/api/leaderboard/map/%s?offset=0&length=15", mapID)
	var lboard Leaderboard
	return getTMObject(UrlLeaderboards, lboard)
}

func getCampaignLeaderboard(leadID string) (Leaderboard, error) {
	// https://trackmania.io/api/leaderboard/NLS-oPaUhtnbfma9lmSHgh8SLFUPJtCK6SZbAxh?offset=0&length=10
	url := fmt.Sprintf("https://trackmania.io/api/leaderboard/%s?offset=0&length=10", leadID)
	var lboard Leaderboard
	return getTMObject(url, lboard)
}

func getCampaignInfo(campID string) (CampaignInfo, error) {
	// campID is something like "60530/45849"
	// https://trackmania.io/api/campaign/60530/45849
	url := fmt.Sprintf("https://trackmania.io/api/campaign/%s", campID)
	var cInfo CampaignInfo
	return getTMObject(url, cInfo)
}

func getTMObject[T interface{}](url string, o T) (T, error) {

	var ret T

	c := http.Client{
		Timeout: time.Second * 5, // Timeout after 5 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ret, err
	}

	res, getErr := c.Do(req)
	if getErr != nil {
		return ret, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return ret, readErr
	}

	jsonErr := json.Unmarshal(body, &ret)
	if jsonErr != nil {
		return ret, jsonErr
	}

	return ret, nil

}

type CampaignInfo struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	Media          string     `json:"media"`
	Creationtime   int        `json:"creationtime"`
	Publishtime    int        `json:"publishtime"`
	Clubid         int        `json:"clubid"`
	Leaderboarduid string     `json:"leaderboarduid"`
	Playlist       []Playlist `json:"playlist"`
	Mediae         Mediae     `json:"mediae"`
	Tracked        bool       `json:"tracked"`
}


type Playlist struct {
	Author          string    `json:"author"`
	Name            string    `json:"name"`
	MapType         string    `json:"mapType"`
	MapStyle        string    `json:"mapStyle"`
	AuthorScore     int       `json:"authorScore"`
	GoldScore       int       `json:"goldScore"`
	SilverScore     int       `json:"silverScore"`
	BronzeScore     int       `json:"bronzeScore"`
	CollectionName  string    `json:"collectionName"`
	Filename        string    `json:"filename"`
	IsPlayable      bool      `json:"isPlayable"`
	MapID           string    `json:"mapId"`
	MapUID          string    `json:"mapUid"`
	Submitter       string    `json:"submitter"`
	Timestamp       time.Time `json:"timestamp"`
	FileURL         string    `json:"fileUrl"`
	ThumbnailURL    string    `json:"thumbnailUrl"`
	Authorplayer    Player    `json:"authorplayer"`
	Submitterplayer Player    `json:"submitterplayer"`
	Exchangeid      int       `json:"exchangeid"`
}
type Mediae struct {
	Buttonbackground     string `json:"buttonbackground"`
	Buttonforeground     string `json:"buttonforeground"`
	Decal                string `json:"decal"`
	Popupbackground      string `json:"popupbackground"`
	Popup                string `json:"popup"`
	Livebuttonbackground string `json:"livebuttonbackground"`
	Livebuttonforeground string `json:"livebuttonforeground"`
}
type Meta struct {
	Twitch  string `json:"twitch"`
	Twitter string `json:"twitter"`
}

// https://trackmania.exchange/api/mappack/get_mappack_tracks/3017
type MapPackMapInfo []struct {
	ShortName            interface{} `json:"ShortName"`
	TrackID              int         `json:"TrackID"`
	UserID               int         `json:"UserID"`
	Username             string      `json:"Username"`
	GbxMapName           string      `json:"GbxMapName"`
	AuthorLogin          string      `json:"AuthorLogin"`
	MapType              string      `json:"MapType"`
	TitlePack            string      `json:"TitlePack"`
	TrackUID             string      `json:"TrackUID"`
	Mood                 string      `json:"Mood"`
	DisplayCost          int         `json:"DisplayCost"`
	ModName              string      `json:"ModName"`
	Lightmap             int         `json:"Lightmap"`
	ExeVersion           string      `json:"ExeVersion"`
	ExeBuild             string      `json:"ExeBuild"`
	AuthorTime           int         `json:"AuthorTime"`
	ParserVersion        int         `json:"ParserVersion"`
	UploadedAt           string      `json:"UploadedAt"`
	UpdatedAt            string      `json:"UpdatedAt"`
	Name                 string      `json:"Name"`
	Tags                 string      `json:"Tags"`
	TypeName             string      `json:"TypeName"`
	StyleName            string      `json:"StyleName"`
	EnvironmentName      string      `json:"EnvironmentName"`
	VehicleName          string      `json:"VehicleName"`
	UnlimiterRequired    bool        `json:"UnlimiterRequired"`
	RouteName            string      `json:"RouteName"`
	LengthName           string      `json:"LengthName"`
	DifficultyName       string      `json:"DifficultyName"`
	Laps                 int         `json:"Laps"`
	ReplayWRID           interface{} `json:"ReplayWRID"`
	ReplayWRTime         interface{} `json:"ReplayWRTime"`
	ReplayWRUserID       interface{} `json:"ReplayWRUserID"`
	ReplayWRUsername     interface{} `json:"ReplayWRUsername"`
	TrackValue           int         `json:"TrackValue"`
	Comments             string      `json:"Comments"`
	MappackID            int         `json:"MappackID"`
	Unlisted             bool        `json:"Unlisted"`
	Unreleased           bool        `json:"Unreleased"`
	Downloadable         bool        `json:"Downloadable"`
	RatingVoteCount      int         `json:"RatingVoteCount"`
	RatingVoteAverage    float64     `json:"RatingVoteAverage"`
	HasScreenshot        bool        `json:"HasScreenshot"`
	HasThumbnail         bool        `json:"HasThumbnail"`
	HasGhostBlocks       bool        `json:"HasGhostBlocks"`
	EmbeddedObjectsCount int         `json:"EmbeddedObjectsCount"`
	EmbeddedItemsSize    int         `json:"EmbeddedItemsSize"`
	AuthorCount          int         `json:"AuthorCount"`
	IsMP4                bool        `json:"IsMP4"`
	SizeWarning          bool        `json:"SizeWarning"`
	AwardCount           int         `json:"AwardCount"`
	CommentCount         int         `json:"CommentCount"`
	ReplayCount          int         `json:"ReplayCount"`
	ImageCount           int         `json:"ImageCount"`
	VideoCount           int         `json:"VideoCount"`
}

// https://trackmania.io/api/leaderboard/map/ZzH2cxK0OTIDunpbtsWTiitfng9?offset=0&length=15
type Leaderboard struct {
	Tops        []LeaderboardPosition `json:"tops"`
	Playercount int                   `json:"playercount"`
}

type Player struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
	ID   string `json:"id"`
	// Zone Zone   `json:"zone"`
	Meta Meta `json:"meta"`
}

type LeaderboardPosition struct {
	Player    Player    `json:"player,omitempty"`
	Position  int       `json:"position"`
	Time      int       `json:"time"`
	Points    int       `json:"points"`
	Filename  string    `json:"filename"`
	Timestamp time.Time `json:"timestamp"`
	URL       string    `json:"url"`
}

//	type Parent struct {
//		Name   string `json:"name"`
//		Flag   string `json:"flag"`
//		Parent Parent `json:"parent"`
//	}
//
//	type Zone struct {
//		Name   string `json:"name"`
//		Flag   string `json:"flag"`
//		Parent Parent `json:"parent"`
//	}
