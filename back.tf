resource "google_service_account" "back" {
  account_id = "back-runner"
}

resource "google_project_iam_member" "back" {
  member  = "serviceAccount:${google_service_account.back.email}"
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
}

resource "google_secret_manager_secret" "api_key" {
  replication {
    automatic = true
  }
  secret_id = "api-key"
}

data "archive_file" "back" {
  output_path = "tmp/back.zip"
  source_dir  = "back"
  type        = "zip"
}

resource "google_storage_bucket" "back" {
  location = var.region
  name     = "${var.project_name}-back"
  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      age = 1
    }
  }
}

resource "google_storage_bucket_object" "back" {
  bucket = google_storage_bucket.back.name
  name   = data.archive_file.back.output_md5
  source = data.archive_file.back.output_path
}

resource "google_cloudfunctions2_function" "back" {
  build_config {
    runtime     = "go120"
    entry_point = "EntryPoint"
    source {
      storage_source {
        bucket = google_storage_bucket.back.name
        object = google_storage_bucket_object.back.name
      }
    }
  }
  location = var.region
  name     = "back"
  service_config {
    min_instance_count = 0
    max_instance_count = 1
    secret_environment_variables {
      key        = "API_KEY"
      project_id = var.project_id
      secret     = google_secret_manager_secret.api_key.secret_id
      version    = "latest"
    }
    service_account_email = google_service_account.back.email
  }
}