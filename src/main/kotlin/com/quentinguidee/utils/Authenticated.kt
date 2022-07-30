package com.quentinguidee.utils

import com.quentinguidee.dao.usersDAO
import com.quentinguidee.models.Session
import com.quentinguidee.models.User
import com.quentinguidee.services.sessionsServices
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.util.pipeline.*
import org.jetbrains.exposed.sql.transactions.transaction

private val sessionKey = io.ktor.util.AttributeKey<Session>("SESSION_KEY")
private val userKey = io.ktor.util.AttributeKey<User>("USER_KEY")

fun Route.authenticated(build: Route.() -> Unit): Route {
    val pipelinePhase = PipelinePhase("Validate")

    val route = createChild(AuthenticatedSelector())

    route.insertPhaseAfter(ApplicationCallPipeline.Plugins, pipelinePhase)
    route.intercept(pipelinePhase) {
        val token = call.request.header(HttpHeaders.Authorization) ?: return@intercept call.respondText(
            "missing authentication token",
            status = HttpStatusCode.BadRequest
        )

        val session = sessionsServices.session(token)
        val user = transaction { usersDAO.get(session.userID) }

        call.attributes.put(sessionKey, session)
        call.attributes.put(userKey, user)
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
