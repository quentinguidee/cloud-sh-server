package com.quentinguidee

import com.quentinguidee.plugins.*
import com.typesafe.config.ConfigFactory
import io.ktor.client.*
import io.ktor.client.engine.cio.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.serialization.kotlinx.json.*
import io.ktor.server.config.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*

val client = HttpClient(CIO) {
    install(ContentNegotiation) {
        json()
    }
}

fun main() {
    embeddedServer(Netty, environment = applicationEngineEnvironment {
        config = HoconApplicationConfig(ConfigFactory.load())

        module {
            configureDatabase()
            configureHTTP()
            configureSerialization()
            configureOAuth()
            configureRouting()
        }

        connector {
            port = 8080
            host = "0.0.0.0"
        }
    }).start(wait = true)
}
