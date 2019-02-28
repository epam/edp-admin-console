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

package models

type App struct {
	Name        string      `json:"name" valid:"Required;Match(/^[a-z]+(-+[a-z0-9]+)*$/)"`
	Strategy    string      `json:"strategy" valid:"Required"`
	Lang        string      `json:"lang" valid:"Required"`
	BuildTool   string      `json:"buildTool" valid:"Required"`
	MultiModule bool        `json:"multiModule,omitempty"`
	Framework   string      `json:"framework" valid:"Required"`
	Repository  *Repository `json:"repository,omitempty"`
	Route       *Route      `json:"route,omitempty"`
	Database    *Database   `json:"database,omitempty"`
	Vcs         *Vcs        `json:"vcs,omitempty"`
}

type Repository struct {
	Url      string `json:"url,omitempty" valid:"Required;Match(/(?:^git|^ssh|^https?|^git@[-\\w.]+):(\\/\\/)?(.*?)(\\.git)(\\/?|\\#[-\\d\\w._]+?)$/)"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type Vcs struct {
	Login    string `json:"login,omitempty" valid:"Required"`
	Password string `json:"password,omitempty" valid:"Required"`
}

type Route struct {
	Site string `json:"site,omitempty" valid:"Required"`
	Path string `json:"path,omitempty" valid:"Required;Match(/^(?:http(s)?:\\/\\/)?[\\w.-]+(?:\\.[\\w\\.-]+)+[\\w\\-\\._~:/?#[\\]@!\\$&'\\(\\)\\*\\+,;=.]+$/)"`
}

type Database struct {
	Kind     string `json:"kind,omitempty" valid:"Required"`
	Version  string `json:"version,omitempty" valid:"Required"`
	Capacity string `json:"capacity,omitempty" valid:"Required"`
	Storage  string `json:"storage,omitempty" valid:"Required"`
}
