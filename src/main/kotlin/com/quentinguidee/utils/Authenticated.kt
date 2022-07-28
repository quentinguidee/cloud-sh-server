package com.quentinguidee.utils

import com.quentinguidee.models.db.Session
import com.quentinguidee.models.db.User
import com.quentinguidee.services.sessionServices
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.util.pipeline.*
import org.jetbrains.exposed.sql.transactions.transaction

private val sessionKey = io.ktor.util.AttributeKey<Session>("SESSION_KEY")
private val userKey = io.ktor.util.AttributeKey<User>("USER_KEY")
private val userIDKey = io.ktor.util.AttributeKey<Int>("USER_ID_KEY")

fun Route.authenticated(build: Route.() -> Unit): Route {
    val pipelinePhase = PipelinePhase("Validate")

    val route = createChild(AuthenticatedSelector())

    route.insertPhaseAfter(ApplicationCallPipeline.Plugins, pipelinePhase)
    route.intercept(pipelinePhase) {
        val token = call.request.header(HttpHeaders.Authorization) ?: return@intercept call.respondText(
            "missing authentication token",
            status = HttpStatusCode.BadRequest
        )

        val session = sessionServices.session(token) ?: return@intercept call.respondText(
            "user session not found",
            status = HttpStatusCode.NotFound
        )

        call.attributes.put(sessionKey, session)
        transaction {
            call.attributes.put(userKey, session.user)
            call.attributes.put(userIDKey, session.user.id.value)
        }
    }
    route.build()

    return route
}

class AuthenticatedSelector : RouteSelector() {
    override fun evaluate(context: RoutingResolveContext, segmentIndex: Int) =
        RouteSelectorEvaluation.Transparent
}

val ApplicationCall.session
    get() = attributes[sessionKey]

val ApplicationCall.user
    get() = attributes[userKey]

val ApplicationCall.userID
    get() = attributes[userIDKey]
