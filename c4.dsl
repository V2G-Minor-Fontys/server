workspace "Automatic Charging Feature" "V2G feature C4 table" {
  !identifiers hierarchical

  model {
    user = person "EV Owner" {
      description "The person who owns the electric vehicle."
    }

    server = softwareSystem "Server" {
      description "Handles client registration, user preferences, and charging schedule logic."

      api = container "API Routes" {
        description "Exposes RESTful endpoints for the mobile app to register, send, and fetch charging-related data."
        technology "Go / REST API"
      }

      scheduler = container "Schedule Generator" {
        description "Responsible for interpreting user preferences and calculating an optimal charging schedule to send to the controller."
        technology "Go"
      }

      mqtt = container "MQTT" {
        description "Acts as a communication endpoint between the client (EV) and the server."
        technology "Go / mqtt"

        server.api -> server.mqtt "Triggers schedule generation"
        server.mqtt -> server.scheduler  "Creates charging schedule"
      }   

      db = container "Data layer" {
        description "Encapsulates all logic for querying and updating the SQL database using sqlc-generated Go code."
        technology "Go / sqlc"

        server.api -> server.db "Requests/Sends data"
        server.scheduler -> server.db "Reads preferences"
      }   
    }

    database = softwareSystem "Database" {
      description "Relational database storing all persistent information such as user profiles, preferences, recurring schedules, and actions."

      db = container "DB" {
        description "PostgreSQL instance backing the data layer for storing normalized entities like users, preferences, and actions."
        technology "PostgreSQL"
      }

      server.db -> database "Reads from and writes to"
    }

    controller = softwareSystem "Client Device (Raspberry Pi)" {
      description "Edge device embedded in the EV setup that communicates with the server to receive up-to-date charging instructions and apply them."
    
      controller -> server.mqtt "Requests charging schedule"
      server.mqtt -> controller "Sends charging schedule"
    }

    app = softwareSystem "Mobile App" {
      tag "Application"
      description "Cross-platform app that provides the EV owner with tools to configure and monitor charging behavior."

      preferences_view = container "Charging Preferences View" {
        technology "Flutter"
        description "Allows users to create, edit, and prioritize charging preferences including times, days, and energy goals."
      }

      user -> app "Uses"
      app.preferences_view -> server.api "Creates and updates user preferences data"
    }
  }

  views {
    systemContext server {
      include *
      include user
      autolayout lr
      title "System Diagram"
    }

    container server {
      include server.api server.scheduler server.mqtt server.db database controller
      title "Container Diagram - Server"
    }

    theme default
  }
}
