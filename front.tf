resource "google_cloud_run_v2_service" "front" {
  name = "front"
  template {}
}