# Artifact store: holds the pipeline's derived output (year_features.json) and
# serves as the export target for BigQuery. The bucket is private — the
# dashboard is served from Vercel, which receives the artifact through the
# deploy pipeline, so nothing here is world-readable.
resource "google_storage_bucket" "artifacts" {
  name                        = var.artifact_bucket_name
  location                    = var.region
  uniform_bucket_level_access = true
  force_destroy               = true
}
