package com.quentinguidee.plugins

import com.quentinguidee.routes.authRoutes
import com.quentinguidee.routes.bucketRoutes
import com.quentinguidee.routes.userRoutes
import com.quentinguidee.utils.authenticated
import io.ktor.server.application.*
import io.ktor.server.routing.*

fun Application.configureRouting() {
    install(IgnoreTrailingSlash)
    routing {
        route("/auth") {
            authRoutes()
        }

        authenticated {
            route("/bucket") {
                bucketRoutes()
            }

            route("/user") {
                userRoutes()
            }
        }
    }
}
