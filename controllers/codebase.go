/*
 * Copyright 2019 EPAM Systems.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controllers

import (
	"edp-admin-console/context"
	"edp-admin-console/models/query"
	"edp-admin-console/service"
	"github.com/astaxie/beego"
	"log"
)

type CodebaseController struct {
	beego.Controller
	CodebaseService  service.CodebaseService
	EDPTenantService service.EDPTenantService
	BranchService    service.CodebaseBranchService
}

func (c *CodebaseController) GetCodebaseOverviewPage() {
	codebaseName := c.GetString(":codebaseName")
	codebase, err := c.CodebaseService.GetCodebaseByName(codebaseName)
	if err != nil {
		c.Abort("500")
		return
	}

	codebase.CodebaseBranch = addCodebaseBranchInProgressIfAny(codebase.CodebaseBranch, c.GetString(paramWaitingForBranch))
	if err != nil {
		c.Abort("500")
		return
	}

	c.Data["EDPVersion"] = context.EDPVersion
	c.Data["Username"] = c.Ctx.Input.Session("username")
	c.Data["Codebase"] = codebase
	c.TplName = "codebase_overview.html"
}

func addCodebaseBranchInProgressIfAny(branches []*query.CodebaseBranch, branchInProgress string) []*query.CodebaseBranch {
	if branchInProgress != "" {
		for _, branch := range branches {
			if branch.Name == branchInProgress {
				return branches
			}
		}

		log.Println("Adding branch " + branchInProgress + " which is going to be created to the list.")
		branch := query.CodebaseBranch{
			Name:   branchInProgress,
			Status: "inactive",
		}
		branches = append(branches, &branch)
	}
	return branches
}
