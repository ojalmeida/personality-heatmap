package models

import "time"

type APIResponse struct {
	Meta struct {
		Status int `json:"status"`
	} `json:"meta"`
	Data struct {
		Results []Result `json:"results,omitempty"`
	} `json:"data"`
}

type User struct {
	ID        string        `json:"_id"`
	Badges    []interface{} `json:"badges"`
	Bio       string        `json:"bio"`
	BirthDate time.Time     `json:"birth_date"`
	Name      string        `json:"name"`
	Photos    []struct {
		ID       string `json:"id"`
		CropInfo struct {
			User struct {
				WidthPct   float64 `json:"width_pct"`
				XOffsetPct float64 `json:"x_offset_pct"`
				HeightPct  float64 `json:"height_pct"`
				YOffsetPct float64 `json:"y_offset_pct"`
			} `json:"user"`
			Algo struct {
				WidthPct   float64 `json:"width_pct"`
				XOffsetPct float64 `json:"x_offset_pct"`
				HeightPct  float64 `json:"height_pct"`
				YOffsetPct float64 `json:"y_offset_pct"`
			} `json:"algo"`
			ProcessedByBullseye bool `json:"processed_by_bullseye"`
			UserCustomized      bool `json:"user_customized"`
			Faces               []struct {
				Algo struct {
					WidthPct   float64 `json:"width_pct"`
					XOffsetPct float64 `json:"x_offset_pct"`
					HeightPct  float64 `json:"height_pct"`
					YOffsetPct float64 `json:"y_offset_pct"`
				} `json:"algo"`
				BoundingBoxPercentage float64 `json:"bounding_box_percentage"`
			} `json:"faces"`
		} `json:"crop_info,omitempty"`
		URL            string `json:"url"`
		ProcessedFiles []struct {
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"processedFiles"`
		ProcessedVideos []interface{} `json:"processedVideos"`
		FileName        string        `json:"fileName"`
		Extension       string        `json:"extension"`
		Assets          []interface{} `json:"assets"`
		MediaType       string        `json:"media_type"`
	} `json:"photos"`
	Gender  int           `json:"gender"`
	Jobs    []interface{} `json:"jobs"`
	Schools []struct {
		Name string `json:"name"`
	} `json:"schools"`
	IsTraveling    bool `json:"is_traveling"`
	RecentlyActive bool `json:"recently_active"`
	OnlineNow      bool `json:"online_now"`
}

type Result struct {
	Type           string `json:"type"`
	User           `json:"user,omitempty"`
	ExperimentInfo struct {
		UserInterests struct {
			SelectedInterests []Interest `json:"selected_interests"`
		} `json:"user_interests"`
	} `json:"experiment_info,omitempty"`
	Spotify struct {
		SpotifyConnected  bool          `json:"spotify_connected"`
		SpotifyTopArtists []interface{} `json:"spotify_top_artists"`
	} `json:"spotify,omitempty"`
	DistanceMi int `json:"distance_mi"`
}

type Profile struct {
	User      User
	Distance  int
	Interests []string
}

type Interest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsCommon bool   `json:"is_common"`
}
