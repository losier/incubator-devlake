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

package tasks

import (
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertApiMergeRequestsMeta = core.SubTaskMeta{
	Name:             "convertApiMergeRequests",
	EntryPoint:       ConvertApiMergeRequests,
	EnabledByDefault: true,
	Description:      "Update domain layer PullRequest according to GitlabMergeRequest",
}

func ConvertApiMergeRequests(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_TABLE)
	db := taskCtx.GetDb()

	domainMrIdGenerator := didgen.NewDomainIdGenerator(&models.GitlabMergeRequest{})
	domainRepoIdGenerator := didgen.NewDomainIdGenerator(&models.GitlabProject{})
	//Find all piplines associated with the current projectid
	cursor, err := db.Model(&models.GitlabMergeRequest{}).Where("project_id=?", data.Options.ProjectId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabMergeRequest{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, error) {
			gitlabMr := inputRow.(*models.GitlabMergeRequest)

			domainPr := &code.PullRequest{
				DomainEntity: domainlayer.DomainEntity{
					Id: domainMrIdGenerator.Generate(gitlabMr.GitlabId),
				},
				BaseRepoId:     domainRepoIdGenerator.Generate(gitlabMr.SourceProjectId),
				HeadRepoId:     domainRepoIdGenerator.Generate(gitlabMr.TargetProjectId),
				Status:         gitlabMr.State,
				Number:         gitlabMr.Iid,
				Title:          gitlabMr.Title,
				Description:    gitlabMr.Description,
				Url:            gitlabMr.WebUrl,
				AuthorName:     gitlabMr.AuthorUsername,
				CreatedDate:    gitlabMr.GitlabCreatedAt,
				MergedDate:     gitlabMr.MergedAt,
				ClosedDate:     gitlabMr.ClosedAt,
				MergeCommitSha: gitlabMr.MergeCommitSha,
				HeadRef:        gitlabMr.SourceBranch,
				BaseRef:        gitlabMr.TargetBranch,
				Component:      gitlabMr.Component,
			}

			return []interface{}{
				domainPr,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
