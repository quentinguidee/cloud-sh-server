package com.quentinguidee

import com.quentinguidee.plugins.configureDatabase
import com.quentinguidee.plugins.configureHTTP
import com.quentinguidee.plugins.configureRouting
import com.quentinguidee.plugins.configureSerialization
import com.typesafe.config.ConfigFactory
import io.ktor.server.config.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*

fun main() {
    embeddedServer(Netty, environment = applicationEngineEnvironment {
        config = HoconApplicationConfig(ConfigFactory.load())

        module {
            configureDatabase()
            configureHTTP()
            configureSerialization()
            configureRouting()
        }

        connector {
            port = 8080
            host = "0.0.0.0"
        }
    }).start(wait = true)
}
