package com.quentinguidee.routes

import com.quentinguidee.services.authService
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.auth.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

fun Route.authRoutes() {
    route("/github") {
        authenticate("oauth-github") {
            get("/login") {}

            get("/callback") {
                val principal: OAuthAccessTokenResponse.OAuth2 = call.principal() ?: return@get call.respondText(
                    "failed to retrieve oauth information",
                    status = HttpStatusCode.InternalServerError
                )

                val githubUserBody = authService.fetchGitHubUser(principal.accessToken)
                val githubUser = authService.githubUser(githubUserBody.login)

                val session = if (githubUser == null) {
                    authService.createAccount(githubUserBody)
                } else {
                    authService.getAccount(githubUser.username)
                }

                call.respond(session.toJSON())
            }
        }
    }
}
