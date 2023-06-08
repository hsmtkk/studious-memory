resource "google_service_account" "front" {
  account_id = "front-runner"
}

resource "google_project_iam_member" "front" {
  member  = "serviceAccount:${google_service_account.front.email}"
  project = var.project_id
  role    = "roles/run.invoker"
}

resource "google_artifact_registry_repository" "registry" {
  format        = "DOCKER"
  repository_id = "registry"
}

resource "google_cloud_run_v2_service" "front" {
  location = var.region
  name     = "front"
  template {
    containers {
      env {
        name  = "BACK_URL"
        value = google_cloudfunctions2_function.back.service_config[0].uri
      }
      image = "path-to-container"
    }
    scaling {
      min_instance_count = 0
      max_instance_count = 1
    }
    service_account = google_service_account.front.email
  }
}