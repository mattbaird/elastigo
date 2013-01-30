elastigo
========

Golang based Elasticsearch client, implements core api for Indexing and searching.

status updates
========================

* 2013-1-26, expansion of search dsl for greater coverage
* 2012-12-30, new bulk indexing and search dsl
* 2012-10-12, early in development, not ready for production yet.

Search Examples
-------------------------

A Faceted, ranged Search using the `Search DSL` :

    out, err := Search("github").Size("1").Facet(
      Facet().Fields("actor").Size("500"),
    ).Query(
      Query().Range(
         Range().Field("created_at").From("2012-12-10T15:00:00-08:00").To("2012-12-10T15:10:00-08:00"),
      ).Search("add"),
    ).Result()
   
A Ranged Search using the `Search DSL` :
   
    out, err := Search("github").Type("Issues").Pretty().Query(
      Query().Range(
         Range().Field("created_at").From("2012-12-10T15:00:00-08:00").To("2012-12-10T15:10:00-08:00"),
      ).Search("add"),
    ).Result()
   
A Simple Search using the `Search DSL` :

    out, err := Search("github").Type("Issues").Size("100).Search("add").Result()


A Direct Search using the api :
   
    qry := map[string]interface{}{
      "query":map[string]interface{}{
         "term":map[string]string{"user:"kimchy"},
      },
    }
    core.SearchRequest(true, "github", "Issues", qry, "")

A Direct Search using the query string Api :
   
    core.SearchUri("github", "Issues", "user:kimchy", "")

A Filtered search `Search DSL` :
   
    out, err := Search("github").Filter(
      Filter().Exists("repository.name"),
    ).Result()


Adding content to Elasticsearch
----------------------------------------------

examples:
    
    // add single go struct entity
    response, _ := core.Index(true, "twitter", "tweet", "1", NewTweet("kimchy", "Search is cool"))

    // you have bytes
    bytesLine, err := json.Marshall(tw)
    response, _ := core.Index(true, "twitter", "tweet", "2", bytesLine)

    // Bulk Indexing 
    core.IndexBulk("twitter", "tweet", "3", &time.Now(), NewTweet("kimchy", "Search is now cooler"))


license
=======
    Copyright 2012 Matthew Baird

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.