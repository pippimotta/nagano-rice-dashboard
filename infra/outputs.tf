output "artifact_bucket" {
  description = "Artifact-store bucket name."
  value       = google_storage_bucket.artifacts.name
}

output "artifact_bucket_uri" {
  description = "gs:// URI of the artifact store."
  value       = "gs://${google_storage_bucket.artifacts.name}"
}
