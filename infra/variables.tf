variable "project_id" {
  description = "GCP project ID that owns the dashboard resources."
  type        = string
}

variable "region" {
  description = "Default region for regional resources (bucket location)."
  type        = string
  default     = "asia-northeast1"
}

variable "artifact_bucket_name" {
  description = "Globally unique name for the private artifact-store bucket."
  type        = string
}
