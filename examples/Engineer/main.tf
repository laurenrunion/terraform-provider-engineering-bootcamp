terraform {
  required_providers {
    engineering-bootcamp = {
      source = "liatr.io/terraform/engineering-bootcamp"
    }
  }
}

provider "engineering-bootcamp" {
  host = "http://localhost:8080"
}

resource "engineering-bootcamp_engineer" "ryan" {
  name  = "Ryan"
  email = "ryan@cats.com"
}

# data source reads the same engineer back through the API
data "engineering-bootcamp_engineer" "ryan" {
  id = engineering-bootcamp_engineer.ryan.id
}

output "ryan_from_datasource" {
  value = data.engineering-bootcamp_engineer.ryan
}
