/*
 * Copyright 2019 Dgraph Labs, Inc. and Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package auth

import (
	"strings"

	dgoapi "github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/dgraph-io/dgraph/gql"
	"github.com/dgraph-io/dgraph/graphql/schema"
)

type BaseProcedure struct {
	sch       *schema.Schema
	authState *schema.AuthState

	queries []*gql.GraphQuery

	re RuleExtractor
}

func (bp *BaseProcedure) Init(sch *schema.Schema, authState *schema.AuthState) {
	bp.sch = sch
	bp.authState = authState
	bp.queries = make([]*gql.GraphQuery, 0)
}

func (bp *BaseProcedure) CollectQueries() []*gql.GraphQuery {
	return bp.queries
}

func (bp *BaseProcedure) getTypeRules(typName string) *schema.RuleNode {
	typRules := (*bp.sch).AuthTypeRules(typName)
	return bp.GetTypeRule(typRules)
}

func (bp *BaseProcedure) getFieldRules(typName, fieldName string) *schema.RuleNode {
	fieldRules := (*bp.sch).AuthFieldRules(typName, fieldName)
	return bp.GetTypeRule(fieldRules)
}

func (bp *BaseProcedure) GetTypeRule(ac *schema.AuthContainer) *schema.RuleNode {
	return bp.re(ac)
}

func (bp *BaseProcedure) SetRuleExtractor(re RuleExtractor) {
	bp.re = re
}

type BaseMutationProcedure struct {
	mutations  []*dgoapi.Mutation
	conditions map[string]bool

	*BaseProcedure
}

func (bmp *BaseMutationProcedure) CollectMutations() []*dgoapi.Mutation {
	return bmp.mutations
}

func (bmp *BaseMutationProcedure) OnMutationCond(cond string) {
	if cond == "" {
		return
	}
	bmp.conditions = make(map[string]bool)

	// trim @if( and )
	cond = cond[4:]
	cond = cond[:len(cond)-1]

	// split by AND
	splits := strings.Split(cond, " AND ")
	for _, split := range splits {
		split = split[7:]
		splitAgain := strings.Split(split, ",")

		varName := splitAgain[0][:len(splitAgain[0])-1]
		varVal := ((splitAgain[1][1] - '0') != 0)
		bmp.conditions[varName] = varVal
	}
}

func QueryRuleExtractor(ac *schema.AuthContainer) *schema.RuleNode {
	if ac == nil {
		return nil
	}
	return ac.Query
}

func AddRuleExtractor(ac *schema.AuthContainer) *schema.RuleNode {
	if ac == nil {
		return nil
	}
	return ac.Add
}

func UpdateRuleExtractor(ac *schema.AuthContainer) *schema.RuleNode {
	if ac == nil {
		return nil
	}
	return ac.Update
}

func DeleteRuleExtractor(ac *schema.AuthContainer) *schema.RuleNode {
	if ac == nil {
		return nil
	}
	return ac.Delete
}