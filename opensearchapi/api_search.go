// SPDX-License-Identifier: Apache-2.0
//
// The OpenSearch Contributors require contributions made to
// this file be licensed under the Apache-2.0 license or a
// compatible open source license.

package opensearchapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/opensearch-project/opensearch-go/v4"
)

// Search executes a /_search request with the optional SearchReq
func (c Client) Search(ctx context.Context, req *SearchReq) (*SearchResp, error) {
	if req == nil {
		req = &SearchReq{}
	}

	var (
		data SearchResp
		err  error
	)
	if data.response, err = c.do(ctx, req, &data); err != nil {
		return &data, err
	}

	return &data, nil
}

// SearchReq represents possible options for the /_search request
type SearchReq struct {
	Indices []string
	Body    io.Reader

	Header http.Header
	Params SearchParams
}

// GetRequest returns the *http.Request that gets executed by the client
func (r SearchReq) GetRequest() (*http.Request, error) {
	var path string
	if len(r.Indices) > 0 {
		path = fmt.Sprintf("/%s/_search", strings.Join(r.Indices, ","))
	} else {
		path = "/_search"
	}

	return opensearch.BuildRequest(
		"POST",
		path,
		r.Body,
		r.Params.get(),
		r.Header,
	)
}

// SearchResp represents the returned struct of the /_search response
type SearchResp struct {
	Took         int                  `json:"took"`
	Timeout      bool                 `json:"timed_out"`
	Shards       ResponseShards       `json:"_shards"`
	Hits         SearchHitsGroup      `json:"hits"`
	Errors       bool                 `json:"errors"`
	Aggregations json.RawMessage      `json:"aggregations"`
	ScrollID     *string              `json:"_scroll_id,omitempty"`
	Suggest      map[string][]Suggest `json:"suggest,omitempty"`
	response     *opensearch.Response
}

// Inspect returns the Inspect type containing the raw *opensearch.Response
func (r SearchResp) Inspect() Inspect {
	return Inspect{Response: r.response}
}

type SearchHitsGroup struct {
	Total 	 HitsTotal 	 `json:"total"`
	MaxScore float32     `json:"max_score"`
	Hits     []SearchHit `json:"hits"`
}

type SearchHitsTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

// SearchHit is a sub type of SearchResp containing information of the search hit with an unparsed Source field
type SearchHit struct {
	Index       string                  `json:"_index"`
	ID          string                  `json:"_id"`
	Routing     string                  `json:"_routing"`
	Score       float32                 `json:"_score"`
	Source      json.RawMessage         `json:"_source"`
	Fields      json.RawMessage         `json:"fields"`
	Type        string                  `json:"_type"` // Deprecated field
	Sort        []any                   `json:"sort"`
	Explanation *DocumentExplainDetails `json:"_explanation"`
	SeqNo       *int                    `json:"_seq_no"`
	PrimaryTerm *int                    `json:"_primary_term"`
}

// Suggest is a sub type of SearchResp containing information of the suggest field
type Suggest struct {
	Text    string `json:"text"`
	Offset  int    `json:"offset"`
	Length  int    `json:"length"`
	Options []struct {
		Text         string  `json:"text"`
		Score        float32 `json:"score"`
		Freq         int     `json:"freq"`
		Highlighted  string  `json:"highlighted"`
		CollateMatch bool    `json:"collate_match"`
	} `json:"options"`
}
