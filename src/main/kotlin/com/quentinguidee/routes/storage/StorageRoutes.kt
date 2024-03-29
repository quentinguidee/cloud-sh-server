package com.quentinguidee.routes.storage

import io.ktor.server.routing.*

fun Route.storageRoutes() {
    route("/storage") {
        binRoutes()
        bucketRoutes()
        nodesRoutes()
        recentRoutes()
    }
}
