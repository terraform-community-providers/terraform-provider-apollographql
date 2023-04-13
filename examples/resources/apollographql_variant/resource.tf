resource "apollographql_variant" "api" {
  name     = "latest"
  graph_id = apollographql_graph.api.id
}
