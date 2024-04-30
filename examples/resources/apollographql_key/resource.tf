resource "apollographql_key" "api" {
  name     = "deploy"
  graph_id = apollographql_graph.api.id
}
