resource "google_cloudbuild_trigger" "cloud_build" {
  filename = "cloudbuild.yaml"
  github {
    owner = "hsmtkk"
    name  = var.project_name
    push {
      branch = "main"
    }
  }
}