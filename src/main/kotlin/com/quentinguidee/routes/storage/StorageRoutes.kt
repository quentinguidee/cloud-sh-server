package com.quentinguidee.routes.storage

import com.quentinguidee.routes.bucketRoutes
import io.ktor.server.routing.*

fun Route.storageRoutes() {
    route("/storage") {
        bucketRoutes()
    }
}
