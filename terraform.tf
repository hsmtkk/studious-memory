terraform {
  backend "gcs" {
    bucket = "studious-memory-tfstate"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}
