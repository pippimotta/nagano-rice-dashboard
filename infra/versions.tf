terraform {
  required_version = ">= 1.5"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }

  # State is local for now (gitignored). A GCS backend can be added later once a
  # dedicated state bucket exists, to share state with CI.
}
