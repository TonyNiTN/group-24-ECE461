terraform {
  required_providers {
    # google = {
    #   source  = "hashicorp/google"
    #   version = "4.51.0"
    # }
    docker = {
      source = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

provider "docker" {
  # Configuration options
  host = "unix:///Users/emile/.docker/run/docker.sock"
}


# resource "google_compute_network" "vpc_network" {
#   name = "terraform-network"
# }

resource "docker_image" "node_image_1" {               
  name = "node_image_1"
  build {  
    context = "../node-app"
  }
}

resource "docker_container" "node_container_1" {   
  # the name of the container
  name = "node_container_1"  
  image = docker_image.node_image_1.image_id
    ports {
        internal = "3000"
        external = "3000"
  }
}

