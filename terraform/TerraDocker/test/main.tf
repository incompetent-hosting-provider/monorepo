terraform {
    required_version = ">= 0.12"
    required_providers {
        random = {
            source = "hashicorp/random"
            version = "2.2.1"
        }
        docker = {
            source = "kreuzwerker/docker"
            version = "2.11.0"
        }
    }
}

provider "docker" {
    host = "unix:///var/run/docker.sock"
}

resource "docker_image" "mysql" {
    name = "mysql:latest"
    keep_locally  = true
}

variable "mysql_default_port" {
    type = number
    default = 3306
    description = "The default port for the mysql container"
}

variable "tfDocker_mysql_name_prefix" {
    type = string
    description = "The prefix for the mysql container name, this variable is set from the environment variables file"
}

variable "tfDocker_uids" {
    type = list(string)
    description = "The list of unique ids for the containers, this variable is set from the environment variables file"

    # This is not possible since the validation does not allow the usage of other variables than the one being validated
    # validation {
    #     condition = length(var.tfDocker_uids) == var.tfDocker_mysql_count
    #     error_message = "The length of the list must be equal to the number of containers you want to create"
    # }
}

variable "mysql_credentials" {
    type = list(string)
    description = "The list of credentials for the mysql container, this variable is set from the creds environment variables file"
}

resource "docker_container" "tfDocker-mysql" {
    count = length(var.tfDocker_uids) < length(var.mysql_credentials) ? length(var.tfDocker_uids) : length(var.mysql_credentials)
    image = docker_image.mysql.latest
    name = "${var.tfDocker_mysql_name_prefix}-${tolist(var.tfDocker_uids)[count.index]}"
    restart = "unless-stopped"

    ports {
        internal = "${var.mysql_default_port}"
        #external = "${var.mysql_default_port+count.index}"
    }

    env = [
        "MYSQL_DATABASE=db",
        "MYSQL_ROOT_PASSWORD=${var.mysql_credentials[count.index]}",
        ]

    volumes {
        container_path = "/var/lib/mysql"
    }

}

output "current_num_mysql_containers" {
    value = length(var.tfDocker_uids) < length(var.mysql_credentials) ? length(var.tfDocker_uids) : length(var.mysql_credentials)
}
