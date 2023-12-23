package utils

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
)

type contest struct {
	Title     string `json:"title"`
	TitleCn   string `json:"titleCn"`
	StartTime int    `json:"startTime"`
}

type userContestRankingInfo struct {
	Data struct {
		UserContestRanking struct {
			AttendedContestsCount   int     `json:"attendedContestsCount"`
			Rating                  float64 `json:"rating"`
			GlobalRanking           int     `json:"globalRanking"`
			LocalRanking            int     `json:"localRanking"`
			GlobalTotalParticipants int     `json:"globalTotalParticipants"`
			LocalTotalParticipants  int     `json:"localTotalParticipants"`
			TopPercentage           float64 `json:"topPercentage"`
		} `json:"userContestRanking"`
		UserContestRankingHistory []struct {
			Attended            bool    `json:"attended"`
			TotalProblems       int     `json:"totalProblems"`
			TrendingDirection   string  `json:"trendingDirection"`
			FinishTimeInSeconds int     `json:"finishTimeInSeconds"`
			Rating              float64 `json:"rating"`
			Score               int     `json:"score"`
			Ranking             int     `json:"ranking"`
			Contest             contest `json:"contest"`
		} `json:"userContestRankingHistory"`
	} `json:"data"`
}

type LcRating struct {
	client http.Client // Define the JSON payload
	req    *http.Request
}

func NewLcRating() (res LcRating, err error) {
	res.client = http.Client{}
	res.req, err = http.NewRequest("POST", "https://leetcode.cn/graphql/noj-go/", nil)
	if err != nil {
		global.GVA_LOG.Warn("请求失败", zap.Error(err))
		return
	}
	res.req.Header.Add("Accept", "*/*")
	res.req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	res.req.Header.Add("Connection", "keep-alive")
	res.req.Header.Add("Content-Type", "application/json")
	res.req.Header.Add("User-Agent", "PostmanRuntime/7.29.2")
	res.req.Header.Add("Content-Length", "1")

	return
}

func (lcRating *LcRating) GetLcRating(LcName string) (Rating float64, err error) {
	Rating = 0
	if LcName == "" {
		return
	}

	// Create a new HTTP client
	jsonStr := `{"query":"\n query userContestRankingInfo($userSlug: String!) {\n userContestRanking(userSlug: $userSlug) {\n attendedContestsCount\n rating\n globalRanking\n localRanking\n globalTotalParticipants\n localTotalParticipants\n topPercentage\n }\n userContestRankingHistory(userSlug: $userSlug) {\n attended\n totalProblems\n trendingDirection\n finishTimeInSeconds\n rating\n score\n ranking\n contest {\n title\n titleCn\n startTime\n }\n }\n}\n ","variables":{"userSlug":"` + LcName + `"},"operationName":"userContestRankingInfo"}`

	lcRating.req.Body = ioutil.NopCloser(strings.NewReader(jsonStr))
	lcRating.req.Header.Set("Content-Length", fmt.Sprintf("%d", len(jsonStr)))

	// Send the request
	resp, err := lcRating.client.Do(lcRating.req)
	if err != nil {
		global.GVA_LOG.Warn("发送请求失败")
		return
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		global.GVA_LOG.Warn("读取响应失败")
		return
	}

	// Close the response body
	resp.Body.Close()

	var rankingInfo userContestRankingInfo
	err = json.Unmarshal(body, &rankingInfo)
	if err != nil {
		global.GVA_LOG.Warn("解析JSON失败")
		return
	}
	return rankingInfo.Data.UserContestRanking.Rating, nil
}
