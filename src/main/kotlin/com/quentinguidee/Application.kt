package com.quentinguidee

import com.quentinguidee.plugins.configureDatabase
import com.quentinguidee.plugins.configureHTTP
import com.quentinguidee.plugins.configureRouting
import com.quentinguidee.plugins.configureSerialization
import io.ktor.server.engine.*
import io.ktor.server.netty.*

fun main() {
    embeddedServer(Netty, port = 8080, host = "0.0.0.0") {
        configureDatabase()
        configureHTTP()
        configureSerialization()
        configureRouting()
    }.start(wait = true)
}
