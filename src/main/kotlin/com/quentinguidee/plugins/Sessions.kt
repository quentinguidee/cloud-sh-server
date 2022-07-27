package com.quentinguidee.plugins

import io.ktor.server.application.*
import io.ktor.server.sessions.*

fun Application.configureSessions() {
    install(Sessions)
}
