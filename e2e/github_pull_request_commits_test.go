/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type GithubPullRequestCommit struct {
	CommitSha string `json:"commit_sha"`
}

func TestGithubPullRequestCommits(t *testing.T) {
	var PullRequestCommits []GithubPullRequestCommit
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT sha FROM github_pull_request_commits where authored_date < '2021-12-25 04:40:11.000'")
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var PullRequestCommit GithubPullRequestCommit
		if err := rows.Scan(&PullRequestCommit.CommitSha); err != nil {
			panic(err)
		}
		PullRequestCommits = append(PullRequestCommits, PullRequestCommit)
	}
	assert.Equal(t, 1505, len(PullRequestCommits))
}
