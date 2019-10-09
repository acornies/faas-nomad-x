job "kafka" {

  type = "service"

  datacenters = ["dc1"]

  group "kafka" {
    count = 1

    task "kafka-brokers" {        
      driver = "docker"
      config {
        image = "confluentinc/cp-enterprise-kafka:5.2.1"
        port_map = {         
          broker = 29092
        }
      }

      env {            
        KAFKA_ZOOKEEPER_CONNECT="$${NOMAD_IP_broker}:2181"
        KAFKA_LISTENER_SECURITY_PROTOCOL_MAP="PLAINTEXT:PLAINTEXT"
        KAFKA_INTER_BROKER_LISTENER_NAME="PLAINTEXT"
        KAFKA_ADVERTISED_LISTENERS="PLAINTEXT://$${NOMAD_IP_broker}:29092"
        KAFKA_BROKER_ID="1"
        KAFKA_BROKER_RACK="r1"
        KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR="1"
        KAFKA_DELETE_TOPIC_ENABLE="true"
        KAFKA_AUTO_CREATE_TOPICS_ENABLE="true"
        KAFKA_LOG4J_ROOT_LOGLEVEL="INFO"
      }
      
      service {
        name = "kafka-broker-$${NOMAD_ALLOC_INDEX}"
        port = "broker"
        check {
          type = "tcp"
          port     = "broker"
          interval = "10s"
          timeout  = "5s"
        }
      }
      resources {
        memory = 512
        network {     
          port "broker" {
            static = 29092
          }
        }
      }
    }
  }

  group "zookeeper" {
    count = 1

    task "zookeeper-svc" {        
      driver = "docker"
      config {
        image = "confluentinc/cp-zookeeper:5.2.1"
        port_map = {         
          client = 2181
        }
      }

      env {            
        ZOOKEEPER_SERVER_ID="1"
        ZOOKEEPER_CLIENT_PORT="2181"
        ZOOKEEPER_TICK_TIME="2000"
      }
      
      service {
        name = "zookeeper-$${NOMAD_ALLOC_INDEX}"
        port = "client"
        check {
          type = "tcp"
          port     = "client"
          interval = "10s"
          timeout  = "5s"
        }
      }
      resources {
        network {     
          port "client" {
            static = 2181
          }
        }
      }
    }
  }

  group "rest-proxy" {
    count = 1

    task "rest-proxy-svc" {        
      driver = "docker"
      config {
        image = "confluentinc/cp-kafka-rest:5.2.1"
        port_map = {         
          http = 8082
        }
      }

      env {            
        KAFKA_REST_ZOOKEEPER_CONNECT="$${NOMAD_IP_http}:2181"
        KAFKA_REST_LISTENERS="http://0.0.0.0:8082"
        KAFKA_REST_SCHEMA_REGISTRY_URL="http://$${NOMAD_IP_http}:8099"
        KAFKA_REST_HOST_NAME="restproxy"
        KAFKA_REST_DEBUG="true"
      }
      
      service {
        name = "kafka-rest-proxy"
        port = "http"
        check {
          type = "tcp"
          port     = "http"
          interval = "10s"
          timeout  = "5s"
        }
      }
      resources {
        network {     
          port "http" {
            static = 8099
          }
        }
      }
    }
  }
}