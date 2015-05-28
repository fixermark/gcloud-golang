// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bigquery

import (
	"reflect"
	"testing"

	bq "google.golang.org/api/bigquery/v2"
)

func defaultQueryJob() *bq.Job {
	return &bq.Job{
		Configuration: &bq.JobConfiguration{
			Query: &bq.JobConfigurationQuery{
				DestinationTable: &bq.TableReference{
					ProjectId: "project-id",
					DatasetId: "dataset-id",
					TableId:   "table-id",
				},
				Query: "query string",
				DefaultDataset: &bq.DatasetReference{
					ProjectId: "def-project-id",
					DatasetId: "def-dataset-id",
				},
			},
		},
	}
}

func TestQuery(t *testing.T) {
	testCases := []struct {
		dst     *Table
		src     *Query
		options []Option
		want    *bq.Job
	}{
		{
			dst:  defaultTable,
			src:  defaultQuery,
			want: defaultQueryJob(),
		},
		{
			dst: defaultTable,
			src: &Query{
				Q: "query string",
			},
			want: func() *bq.Job {
				j := defaultQueryJob()
				j.Configuration.Query.DefaultDataset = nil
				return j
			}(),
		},
		{
			dst: &Table{},
			src: defaultQuery,
			want: func() *bq.Job {
				j := defaultQueryJob()
				j.Configuration.Query.DestinationTable = nil
				return j
			}(),
		},
		{
			dst: &Table{
				ProjectID:         "project-id",
				DatasetID:         "dataset-id",
				TableID:           "table-id",
				CreateDisposition: "CREATE_NEVER",
				WriteDisposition:  "WRITE_TRUNCATE",
			},
			src: defaultQuery,
			want: func() *bq.Job {
				j := defaultQueryJob()
				j.Configuration.Query.WriteDisposition = "WRITE_TRUNCATE"
				j.Configuration.Query.CreateDisposition = "CREATE_NEVER"
				return j
			}(),
		},
	}

	for _, tc := range testCases {
		c := &testClient{}
		if _, err := query(tc.dst, tc.src, c, tc.options); err != nil {
			t.Errorf("err calling query: %v", err)
			continue
		}
		if !reflect.DeepEqual(c.Job, tc.want) {
			t.Errorf("insertJob got:\n%v\nwant:\n%v", c.Job, tc.want)
		}
	}
}