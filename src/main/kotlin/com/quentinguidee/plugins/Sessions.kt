package com.quentinguidee.plugins

import com.quentinguidee.models.UserSession
import io.ktor.server.application.*
import io.ktor.server.sessions.*
import java.io.File

fun Application.configureSessions() {
    install(Sessions) {
        cookie<UserSession>("user_session", directorySessionStorage(File("build/sessions")))
    }
}
